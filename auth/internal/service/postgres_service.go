package service

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/models"
	"context"

	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

var _ Service = (*DBService)(nil)

func NewService(cfg *config.Config) (Service, error) {
	db, err := database.NewDatabase(cfg, &models.User{})
	if err != nil {
		return nil, err
	}
	return &DBService{
		db: db,
	}, nil
}

func (p *DBService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, gorm.ErrInvalidData
	}

	hashedPassword, err := user.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := p.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *DBService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return gorm.ErrRecordNotFound
	}

	if err := p.db.WithContext(ctx).Delete(&models.User{ID: id}).Error; err != nil {
		return nil
	}

	return nil
}
func (p *DBService) Read(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user models.User

	if err := p.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *DBService) Update(ctx context.Context, user *models.User) error {
	if user == nil || user.ID <= 0 {
		return gorm.ErrInvalidData
	}
	updates := make(map[string]interface{})

	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Password != "" {
		hashedPassword, err := user.HashPassword(user.Password)
		if err != nil {
			return err
		}
		updates["password"] = hashedPassword
	}

	if len(updates) == 0 {
		return nil
	}

	result := p.db.WithContext(ctx).Model(&models.User{ID: user.ID}).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (p *DBService) ReadByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, gorm.ErrInvalidData
	}

	var user models.User

	if err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *DBService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, gorm.ErrInvalidData
	}

	var user models.User

	if err := p.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	if !user.CheckPassword(password, user.Password) {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

func (p *DBService) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
