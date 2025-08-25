package controller

import (
	"net/http"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/mapper"
	"github.com/ducklawrence05/go-test-backend-api/internal/controller/http/v1/request"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	auth    user.UserAuthManager
	profile user.UserProfileManager
}

func NewUserController(
	auth user.UserAuthManager,
	profile user.UserProfileManager,
) *UserController {
	return &UserController{
		auth:    auth,
		profile: profile,
	}
}

func (uc *UserController) SendRegistrationOTP(c *gin.Context) {
	var req request.SendRegistrationOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	if err := uc.auth.SendRegistrationOTP(ctx, req.Email); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to get OTP",
	})
}

func (uc *UserController) VerifyRegistrationOTP(c *gin.Context) {
	var req request.VerifyRegistrationOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	emailVerifiedToken, err := uc.auth.VerifyRegistrationOTP(ctx, req.Email, req.OTP)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "Verify email success",
		"email_verified_token": emailVerifiedToken,
	})
}

func (uc *UserController) CompleteRegistration(c *gin.Context) {
	var req request.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing email in email verified token"})
		return
	}
	vo := user.CreateUserVO{
		Email:     email.(string),
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	accessToken, refreshToken, err := uc.auth.CompleteRegistration(ctx, vo)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register success",
		"token": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (uc *UserController) Login(c *gin.Context) {
	var req request.LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	vo := user.LoginUserVO{
		EmailOrUsername: req.UserName,
		Password:        req.Password,
	}

	ctx := c.Request.Context()

	accessToken, refreshToken, err := uc.auth.Login(ctx, vo)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"token": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func (uc *UserController) Logout(c *gin.Context) {
	var req request.LogoutUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	vo := user.LogoutUserVO{
		UserID:       userID.(uuid.UUID),
		RefreshToken: req.RefreshToken,
	}

	ctx := c.Request.Context()

	err := uc.auth.Logout(ctx, vo)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (uc *UserController) GetMe(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	ctx := c.Request.Context()

	user, err := uc.profile.GetMe(ctx, userID.(uuid.UUID))
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, mapper.ToUserInfoResponse(user))
}

func (uc *UserController) UpdateMe(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	var req request.UpdateMeUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	vo := user.UpdateMeVO{
		UserID:    userID.(uuid.UUID),
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	ctx := c.Request.Context()

	user, err := uc.profile.UpdateMe(ctx, vo)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, mapper.ToUserInfoResponse(user))
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	var req request.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	vo := user.ChangePasswordVO{
		UserID:      userID.(uuid.UUID),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	ctx := c.Request.Context()

	if err := uc.profile.ChangePassword(ctx, vo); err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "change password success"})
}

func (uc *UserController) DeleteMe(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	ctx := c.Request.Context()

	err := uc.profile.DeleteMe(ctx, userID.(uuid.UUID))
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
}

func (uc *UserController) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validation.TranslateValidationError(err),
		})
		return
	}

	ctx := c.Request.Context()

	accessToken, refreshToken, err := uc.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		errorcode.JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "refresh token success",
		"token": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}
