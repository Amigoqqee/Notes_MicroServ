package service

import (
	"auth/internal/models"
	"context"
)

type Service interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Read(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	Authenticate(ctx context.Context, username, password string) (*models.User, error)
	ReadByUsername(ctx context.Context, username string) (*models.User, error)
	Close() error
}
