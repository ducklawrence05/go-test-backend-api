package handler

import (
	"net/http"

	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/payload"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"

	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils/mapper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) Register(c *gin.Context) {
	var payload payload.RegisterUserPayLoad
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.TranslateValidationError(err),
		})
		return
	}

	vo := user.CreateUserVO{
		UserName:  payload.UserName,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
	}

	if err := uh.userService.Register(vo); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "register success"})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var payload payload.LoginUserPayLoad
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.TranslateValidationError(err),
		})
		return
	}

	vo := user.LoginUserVO{
		UserName: payload.UserName,
		Password: payload.Password,
	}

	accessToken, refreshToken, err := uh.userService.Login(vo)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
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

func (uh *UserHandler) Logout(c *gin.Context) {
	var payload payload.LogoutUserPayLoad
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.TranslateValidationError(err),
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
		RefreshToken: payload.RefreshToken,
	}

	err := uh.userService.Logout(vo)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (uh *UserHandler) GetMe(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	user, err := uh.userService.GetMe(userID.(uuid.UUID))
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.IndentedJSON(http.StatusOK, mapper.ToUserInfoResponse(user))
}

func (uh *UserHandler) UpdateMe(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	var payload payload.UpdateMeUserPayLoad
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.TranslateValidationError(err),
		})
		return
	}

	vo := user.UpdateMeVO{
		UserID:    userID.(uuid.UUID),
		UserName:  payload.UserName,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}

	user, err := uh.userService.UpdateMe(vo)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.IndentedJSON(http.StatusOK, mapper.ToUserInfoResponse(user))
}

func (uh *UserHandler) ChangePassword(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	var payload payload.ChangePasswordPayLoad
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.TranslateValidationError(err),
		})
		return
	}

	vo := user.ChangePasswordVO{
		UserID:      userID.(uuid.UUID),
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}

	if err := uh.userService.ChangePassword(vo); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "change password success"})
}

func (uh *UserHandler) DeleteMe(c *gin.Context) {
	// get userID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing userID in access token"})
		return
	}

	err := uh.userService.DeleteMe(userID.(uuid.UUID))
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
}

func (uh *UserHandler) RefreshToken(c *gin.Context) {
	var payload payload.RefreshTokenPayLoad
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.TranslateValidationError(err),
		})
		return
	}

	accessToken, refreshToken, err := uh.userService.RefreshToken(payload.RefreshToken)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{
			"error": err.Msg,
		})
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
