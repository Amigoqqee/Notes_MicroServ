package service

import (
	"context"
	"notes/internal/models"
)

type Service interface {
	Close() error
	Create(ctx context.Context, note models.Note) (*models.Note, error)
	GetByID(ctx context.Context, id string) (*models.Note, error)
	GetAll(ctx context.Context, authorId int) ([]models.Note, error)
	Update(ctx context.Context, note models.Note) (*models.Note, error)
	Delete(ctx context.Context, id string) error
}
