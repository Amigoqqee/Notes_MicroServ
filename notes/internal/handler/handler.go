package handler

import (
	"context"
	jwtmanager "jwt_manager"
	"net/http"
	"notes/internal/config"
	"notes/internal/errors"
	"notes/internal/models"
	"notes/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg        *config.Config
	jwtManager *jwtmanager.JWTManager
	service    service.Service
}

func NewHandler(cfg *config.Config, service service.Service) *Handler {
	jwtConfig := jwtmanager.JWTConfig{
		SecretKey:              cfg.JWTSecretKey,
		AccessTokenExpiration:  24,
		RefreshTokenExpiration: 168,
	}
	jwtManager := jwtmanager.NewJWTManager(jwtConfig)

	return &Handler{
		cfg:        cfg,
		jwtManager: jwtManager,
		service:    service,
	}
}

func (h *Handler) CreateNote(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}

	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	note.AuthorID = authorID

	ctx := context.Background()
	createdNote, err := h.service.Create(ctx, note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgNoteCreation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": errors.MsgNoteCreated,
		"note":    createdNote,
	})
}

func (h *Handler) GetNoteByID(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.MsgInvalidNoteID,
		})
		return
	}

	ctx := context.Background()
	note, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   errors.MsgNoteNotFound,
			"details": err.Error(),
		})
		return
	}

	if note.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgNoteFound,
		"note":    note,
	})
}

func (h *Handler) UpdateNote(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.MsgInvalidNoteID,
		})
		return
	}

	ctx := context.Background()
	existingNote, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   errors.MsgNoteNotFound,
			"details": err.Error(),
		})
		return
	}

	if existingNote.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	note.ID = id
	note.AuthorID = authorID

	updatedNote, err := h.service.Update(ctx, note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgNoteUpdate,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgNoteUpdated,
		"note":    updatedNote,
	})
}

func (h *Handler) DeleteNote(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.MsgInvalidNoteID,
		})
		return
	}

	ctx := context.Background()
	existingNote, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   errors.MsgNoteNotFound,
			"details": err.Error(),
		})
		return
	}

	if existingNote.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgNoteDeletion,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": errors.MsgNoteDeleted,
	})
}

func (h *Handler) GetAllNotes(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   errors.MsgMissingUserID,
			"details": err.Error(),
		})
		return
	}

	ctx := context.Background()
	notes, err := h.service.GetAll(ctx, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   errors.MsgNotesFound,
		"notes":     notes,
		"count":     len(notes),
		"author_id": authorID,
	})
}

func (h *Handler) GetJWTMiddleware() gin.HandlerFunc {
	return h.jwtManager.JWTInterceptor()
}

func (h *Handler) extractAuthorID(c *gin.Context) (int, error) {
	return jwtmanager.GetCurrentUserID(c)
}
