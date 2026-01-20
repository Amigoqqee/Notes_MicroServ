package database

import (
	"context"
	"fmt"
	"notes/internal/config"
	"notes/internal/errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase(cfg *config.Config) (*mongo.Client, error) {
	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("%w", errors.ErrEmptyDSN)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DBTimeout)*time.Second)
	defer cancel()

	db, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBDSN))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	return db, nil
}

func CloseDB(db *mongo.Client, cgf *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cgf.Timeout)*time.Second)
	defer cancel()
	if db == nil {
		return fmt.Errorf("%w", errors.ErrDatabaseNotInit)
	}

	return db.Disconnect(ctx)
}
