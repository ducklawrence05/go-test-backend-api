package handler

import (
	"net/http"

	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/payload"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/user"

	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/mapper"
	"github.com/ducklawrence05/go-test-backend-api/pkg/utils"
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

	ctx := c.Request.Context()

	accessToken, refreshToken, err := uh.userService.Register(ctx, vo)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"token": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
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

	ctx := c.Request.Context()

	accessToken, refreshToken, err := uh.userService.Login(ctx, vo)
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

	ctx := c.Request.Context()

	err := uh.userService.Logout(ctx, vo)
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

	ctx := c.Request.Context()

	user, err := uh.userService.GetMe(ctx, userID.(uuid.UUID))
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

	ctx := c.Request.Context()

	user, err := uh.userService.UpdateMe(ctx, vo)
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

	ctx := c.Request.Context()

	if err := uh.userService.ChangePassword(ctx, vo); err != nil {
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

	ctx := c.Request.Context()

	err := uh.userService.DeleteMe(ctx, userID.(uuid.UUID))
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

	ctx := c.Request.Context()

	accessToken, refreshToken, err := uh.userService.RefreshToken(ctx, payload.RefreshToken)
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
