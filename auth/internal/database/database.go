package database

import (
	"auth/internal/config"
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config, models ...any) (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	time.Sleep(1 * time.Second)

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", "ошибка подключения к базе данных", err)
	}

	errMigration := runMigrations(db, models...)
	if errMigration != nil {
		return nil, fmt.Errorf("%s: %v", "ошибка миграции базы данных", errMigration)
	}

	return db.WithContext(ctx), nil
}
func runMigrations(db *gorm.DB, models ...any) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("ошибка миграции модели %T: %w", model, err)
		}
	}

	fmt.Println("Все миграции успешнo выполнены")
	return nil
}
