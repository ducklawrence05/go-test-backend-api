package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/repository"
	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userService struct {
	config           *setting.Config
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewService(
	config *setting.Config,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) Service {
	return &userService{
		config:           config,
		userRepo:         userRepo,
		roleRepo:         roleRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (us *userService) createAccessToken(user_id uuid.UUID) (string, error) {
	return utils.CreateJWT(
		[]byte(us.config.JWT.AccessTokenKey),
		user_id,
		us.config.JWT.AccessTokenExpiresIn,
	)
}

func (us *userService) createRefreshToken(user_id uuid.UUID) (string, error) {
	return utils.CreateJWT(
		[]byte(us.config.JWT.RefreshTokenKey),
		user_id,
		us.config.JWT.RefreshTokenExpiresIn,
	)
}

func (us *userService) checkRefreshToken(token string, userID uuid.UUID) (*entities.RefreshToken, error) {
	refreshToken, err := us.refreshTokenRepo.GetByTokenAndUserID(token, userID)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

// GenerateTokens concurrently creates access token and refresh token
func (us *userService) generateTokens(userID uuid.UUID) (accessToken string, refreshToken string, myErr *utils.MyError) {
	type tokenResult struct {
		Token string
		Err   error
	}

	acCh := make(chan tokenResult)
	rtCh := make(chan tokenResult)

	// generate access token
	go func() {
		at, err := us.createAccessToken(userID)
		acCh <- tokenResult{at, err}
	}()

	// generate refresh token
	go func() {
		rt, err := us.createRefreshToken(userID)
		rtCh <- tokenResult{rt, err}
	}()

	// receive results
	acRes := <-acCh
	rtRes := <-rtCh

	if acRes.Err != nil {
		return "", "", &utils.MyError{
			Msg:        acRes.Err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	if rtRes.Err != nil {
		return "", "", &utils.MyError{
			Msg:        rtRes.Err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return acRes.Token, rtRes.Token, nil
}

// Register implements IUserService.
func (us *userService) Register(vo CreateUserVO) *utils.MyError {
	// check if user exists
	if _, err := us.userRepo.GetByUsername(vo.UserName); err == nil {
		return &utils.MyError{
			Msg:        "this username is already exists",
			StatusCode: http.StatusBadRequest,
		}
	}

	// hash pass
	hashedPassword, err := utils.HashPassword(vo.Password)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// get default role
	defaultRole, err := us.roleRepo.GetByName("user")
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = us.userRepo.Create(&entities.User{
		ID:        uuid.New(),
		UserName:  vo.UserName,
		FirstName: vo.FirstName,
		LastName:  vo.LastName,
		Password:  hashedPassword,
		IsActive:  true,
		CreatedAt: time.Now(),
		RoleID:    defaultRole.ID,
	})

	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *userService) Login(vo LoginUserVO) (string, string, *utils.MyError) {
	// get user from db
	user, err := us.userRepo.GetByUsername(vo.UserName)
	if err != nil || !utils.ComparePasswords(user.Password, []byte(vo.Password)) {
		return "", "", &utils.MyError{
			Msg:        "invalid username or password",
			StatusCode: http.StatusBadRequest,
		}
	}

	// gene ac and rt
	accessToken, refreshToken, myErr := us.generateTokens(user.ID)
	if myErr != nil {
		return "", "", myErr
	}

	// decode rt to get exp and iat
	claims, err := utils.ValidateToken(
		[]byte(us.config.JWT.RefreshTokenKey),
		refreshToken,
	)
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// insert rt to into db
	err = us.refreshTokenRepo.Create(&entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		CreatedAt: time.Now(),
		Revoked:   false,
	})
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return accessToken, refreshToken, nil
}

func (us *userService) Logout(vo LogoutUserVO) *utils.MyError {
	// decode rt
	claims, err := utils.ValidateToken(
		[]byte(us.config.JWT.RefreshTokenKey),
		vo.RefreshToken,
	)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	if claims.UserID != vo.UserID {
		return &utils.MyError{
			Msg:        "refresh token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	}

	_, err = us.checkRefreshToken(vo.RefreshToken, vo.UserID)
	if err != nil {
		return &utils.MyError{
			Msg:        "refresh token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	}

	err = us.refreshTokenRepo.Revoke(vo.UserID, vo.RefreshToken)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *userService) GetMe(userID uuid.UUID) (*entities.User, *utils.MyError) {
	user, err := us.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (us *userService) UpdateMe(vo UpdateMeVO) (*entities.User, *utils.MyError) {
	// get user
	user, err := us.userRepo.GetByID(vo.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// check username unique
	if vo.UserName != "" {
		taken, err := us.userRepo.IsUserNameTaken(vo.UserName, user.ID)
		if err != nil {
			return nil, &utils.MyError{
				Msg:        err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}
		if taken {
			return nil, &utils.MyError{
				Msg:        "username already taken",
				StatusCode: http.StatusUnprocessableEntity,
			}
		}
		user.UserName = vo.UserName
	}

	// other fields
	if vo.FirstName != "" {
		user.FirstName = vo.FirstName
	}
	if vo.LastName != "" {
		user.LastName = vo.LastName
	}

	if err := us.userRepo.Update(user, map[string]any{
		"user_name":  user.UserName,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}); err != nil {
		return nil, &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return user, nil
}

func (us *userService) ChangePassword(vo ChangePasswordVO) *utils.MyError {
	// get user
	user, err := us.userRepo.GetByID(vo.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// check old password
	if !utils.ComparePasswords(user.Password, []byte(vo.OldPassword)) {
		return &utils.MyError{
			Msg:        "invalid password",
			StatusCode: http.StatusBadRequest,
		}
	}

	// hash password
	hashedPassword, err := utils.HashPassword(vo.NewPassword)
	if err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// change password
	if err := us.userRepo.Update(user, map[string]any{
		"password": hashedPassword,
	}); err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *userService) DeleteMe(userID uuid.UUID) *utils.MyError {
	// check user exists
	if _, err := us.userRepo.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.MyError{
				Msg:        "user not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// hard delete all rt
	if err := us.refreshTokenRepo.DeleteByUserID(userID); err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// hard delete user
	if err := us.userRepo.DeleteByID(userID); err != nil {
		return &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func (us *userService) RefreshToken(refreshToken string) (string, string, *utils.MyError) {
	// valid token
	claims, err := utils.ValidateToken([]byte(us.config.JWT.RefreshTokenKey), refreshToken)
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        "refresh token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	}

	// check token in db
	if _, err := us.checkRefreshToken(refreshToken, claims.UserID); err != nil {
		return "", "", &utils.MyError{
			Msg:        "refresh token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	}

	// gene ac and rt
	accessToken, newRefreshToken, myErr := us.generateTokens(claims.UserID)
	if myErr != nil {
		return "", "", myErr
	}

	// decode rt to get exp and iat
	newClaims, err := utils.ValidateToken(
		[]byte(us.config.JWT.RefreshTokenKey),
		newRefreshToken,
	)
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// insert rt to into db
	err = us.refreshTokenRepo.Create(&entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    newClaims.UserID,
		Token:     newRefreshToken,
		IssuedAt:  newClaims.IssuedAt.Time,
		ExpiresAt: newClaims.ExpiresAt.Time,
		CreatedAt: time.Now(),
		Revoked:   false,
	})
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// revoke old rt
	err = us.refreshTokenRepo.Revoke(newClaims.UserID, refreshToken)
	if err != nil {
		return "", "", &utils.MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return accessToken, newRefreshToken, nil
}
