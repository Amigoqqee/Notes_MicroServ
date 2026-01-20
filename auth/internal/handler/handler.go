package handler

import (
	"auth/internal/config"
	"auth/internal/errors"
	"auth/internal/models"
	"auth/internal/service"
	"context"
	jwtmanager "jwt_manager"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    service.Service
	jwtManager *jwtmanager.JWTManager
	cfg        *config.Config
}

func NewHandler(service service.Service, cfg *config.Config) *Handler {
	jwtConfig := jwtmanager.JWTConfig{
		SecretKey:              cfg.JWTSecretKey,
		AccessTokenExpiration:  cfg.AccessTokenExpiration,
		RefreshTokenExpiration: cfg.RefreshTokenExpiration,
	}
	jwtManager := jwtmanager.NewJWTManager(jwtConfig)

	return &Handler{
		service:    service,
		jwtManager: jwtManager,
		cfg:        cfg,
	}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(400, gin.H{
			"error": errors.ErrInvalidUserData,
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	createdUser, err := h.service.Create(ctx, &user)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.ErrUserCreation,
			"details": err.Error(),
		})
		return
	}

	createdUser.Password = ""

	c.JSON(201, gin.H{
		"message": errors.MsgUserRegistered,
		"user":    createdUser,
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	user, err := h.service.Authenticate(ctx, loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgInvalidCredentials,
		})
		return
	}

	user.Password = ""

	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":       errors.MsgLoginSuccess,
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.ErrAuthRequired,
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	user, err := h.service.Read(ctx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserIdNotFound,
		})
		return
	}

	user.Password = ""

	c.JSON(200, gin.H{
		"user": user,
	})
}
func (h *Handler) GetCurrentUserID(c *gin.Context) (int, error) {
	return jwtmanager.GetCurrentUserID(c)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	var updateData models.User

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	updateData.ID = userID

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	err = h.service.Update(ctx, &updateData)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	readCtx, readCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer readCancel()

	updatedUser, err := h.service.Read(readCtx, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": errors.MsgDatabaseOperation,
		})
		return
	}

	updatedUser.Password = ""

	c.JSON(200, gin.H{
		"message": errors.MsgUserUpdated,
		"user":    updatedUser,
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	checkCtx, checkCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer checkCancel()

	_, err = h.service.Read(checkCtx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	deleteCtx, deleteCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer deleteCancel()

	err = h.service.Delete(deleteCtx, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": errors.MsgUserDeleted,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	userID, err := h.jwtManager.ValidateRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgRefreshToken,
		})
		return
	}

	user, err := h.service.Read(c.Request.Context(), userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":       errors.MsgTokensRefreshed,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func (h *Handler) ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.ErrMissingAuthHeader
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.ErrInvalidAuthFormat
	}

	return authHeader[len(bearerPrefix):], nil
}

func (h *Handler) ValidateAccessToken(tokenString string) (int, error) {
	return h.jwtManager.ValidateAccessToken(tokenString)
}

func (h *Handler) RequireAuth() gin.HandlerFunc {
	return h.jwtManager.JWTInterceptor()
}
