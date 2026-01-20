package service

import (
	"context"
	"encoding/json"
	"fmt"
	"notes/internal/caching"
	"notes/internal/config"
	"notes/internal/database"
	"notes/internal/errors"
	"notes/internal/models"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoService struct {
	db         *mongo.Client
	collection *mongo.Collection
	caching    *redis.Client
}

var _ Service = (*MongoService)(nil)

func NewService(cfg *config.Config) (Service, error) {
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseConnection, err)
	}

	cache, err := caching.NewCaching(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrCacheConnection, err)
	}

	collection := db.Database(cfg.DB_NAME).Collection(cfg.DB_COLLECTION)

	return &MongoService{
		db:         db,
		collection: collection,
		caching:    cache,
	}, nil
}
func (m *MongoService) Create(ctx context.Context, note models.Note) (*models.Note, error) {
	result, err := m.collection.InsertOne(ctx, bson.M{
		"name":      note.Name,
		"content":   note.Content,
		"author_id": note.AuthorID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteCreation, err)
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	note.ID = insertedID.Hex()

	m.invalidateAuthorCache(note.AuthorID)

	return &note, nil
}

func (m *MongoService) GetByID(ctx context.Context, id string) (*models.Note, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	var note models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
		}
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	note.ID = objectID.Hex()

	return &note, nil
}

func (m *MongoService) GetAll(ctx context.Context, authorId int) ([]models.Note, error) {
	if cachedNotes, found := m.getCachedNotes(authorId); found {
		fmt.Println("Заметки получены из кэша")
		return cachedNotes, nil
	}

	filter := bson.M{"author_id": authorId}

	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrDecodeNote, err)
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrIterationNotes, err)
	}

	m.cacheNotes(authorId, notes)

	return notes, nil
}

func (m *MongoService) Update(ctx context.Context, note models.Note) (*models.Note, error) {
	objectID, err := primitive.ObjectIDFromHex(note.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	var existingNote models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingNote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, note.ID)
		}
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":    note.Name,
			"content": note.Content,
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteUpdate, err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, note.ID)
	}

	m.invalidateAuthorCache(existingNote.AuthorID)

	note.AuthorID = existingNote.AuthorID

	return &note, nil
}

func (m *MongoService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}

	var existingNote models.Note
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingNote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
		}
		return fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNoteDeletion, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: заметка с ID %s не найдена", errors.ErrNoteNotFound, id)
	}

	m.invalidateAuthorCache(existingNote.AuthorID)

	return nil
}

func (m *MongoService) Close() error {
	if m.caching != nil {
		if err := m.caching.Close(); err != nil {
			return fmt.Errorf("%w: %v", errors.ErrCacheClose, err)
		}
	}
	return database.CloseDB(m.db, &config.Config{Timeout: 10})
}

func (m *MongoService) getCachedNotes(authorID int) ([]models.Note, bool) {
	cacheKey := m.getCacheKey(authorID)
	cachedData, err := m.caching.Get(cacheKey).Result()
	if err != nil {
		fmt.Printf("Ошибка при получении кэша для автора с ID %d: %v\n", authorID, err)
		return nil, false
	}
	var cachedNotes []models.Note
	if err := json.Unmarshal([]byte(cachedData), &cachedNotes); err != nil {
		fmt.Printf("Ошибка при разборе кэша для автора с ID %d: %v\n", authorID, err)
		return nil, false
	}
	fmt.Println("Заметки для автора с ID", authorID, "успешно получены из кэша")
	return cachedNotes, true
}

func (m *MongoService) getCacheKey(authorID int) string {
	return fmt.Sprintf("notes:author:%d", authorID)
}

func (m *MongoService) invalidateAuthorCache(authorID int) {
	if m.caching != nil {
		cacheKey := m.getCacheKey(authorID)
		m.caching.Del(cacheKey)
		fmt.Println("Кэш для автора с ID", authorID, "был успешно инвалидирован")
	}
	fmt.Println("Кэш для автора с ID", authorID, "не найден или не был инвалидирован")
}

func (m *MongoService) cacheNotes(authorID int, notes []models.Note) {
	if m.caching != nil {
		cacheKey := m.getCacheKey(authorID)
		notesJSON, err := json.Marshal(notes)
		if err == nil {
			m.caching.Set(cacheKey, notesJSON, 100*time.Minute)
			fmt.Println("Заметки для автора с ID", authorID, "успешно сохранены в кэш")
		}
	}
}
