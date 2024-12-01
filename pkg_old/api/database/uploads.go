package database

import (
	"context"
	"fmt"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

// AddMedia adds new media to database
func AddMedia(id string, username string, path string, size int64) error {
	media := &models.Media{
		ID:       id,
		Username: username,
		Path:     path,
		Size:     size,
	}

	_, err := GetCollection("media").InsertOne(context.TODO(), media)
	if err != nil {
		return fmt.Errorf("failed to insert media: %w", err)
	}
	return nil
}

// GetMedia retrieves media by ID
func GetMedia(id string) (*models.Media, error) {
	var media models.Media
	err := GetCollection("media").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&media)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch media: %w", err)
	}
	return &media, nil
}

// DeleteMedia removes media from database
func DeleteMedia(id string) error {
	result, err := GetCollection("media").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("media not found: %s", id)
	}
	return nil
}
