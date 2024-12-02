package database

import (
	"context"
	"fmt"
	"os"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a reference to on-disk media to the database, allowing you to specify the
id of this media... useful for any default media.
*/
func AddDefaultMedia(id string, path string, username string) (string, error) {
	// Creating the new media reference
	media := &models.Media{
		ID:       id,
		Path:     path,
		AuthorID: username,
	}

	// Trying to add media to database assuming its path isn't already taken
	mediaCollection := GetCollection("media")
	err := mediaCollection.FindOne(context.TODO(), gin.H{"Path": path}).Decode(&media)
	if err == nil {
		return "", fmt.Errorf("media path already exists")
	}
	_, err = mediaCollection.InsertOne(context.TODO(), media)
	if err != nil {
		return "", err
	}

	return media.ID, nil
}

/*
Adds a reference to on-disk media to the database.
*/
func AddMedia(path string, username string) (string, error) {
	return AddDefaultMedia(uuid.New().String(), path, username)
}

/*
Removes media model from the database.
*/
func RemoveMedia(id string) error {
	mediaCollection := GetCollection("media")
	_, err := mediaCollection.DeleteOne(context.TODO(), gin.H{"_id": id})
	return err
}

/*
Removes media model from the database and from disk.
*/
func RemoveMediaFromDisk(id string) error {
	// Grabbing media model and verifying it exists
	media, err := GetMedia(id)
	if err != nil {
		return err
	}

	// Remove media from database
	if err = RemoveMedia(id); err != nil {
		return err
	}

	// Deleting media from disk
	if err = os.Remove("./" + media.Path); err != nil {
		return err
	}

	return nil
}

/*
Gets media model from the database.
*/
func GetMedia(id string) (*models.Media, error) {
	var media *models.Media
	err := GetCollection("media").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&media)
	return media, err
}

/*
Grabs path of on-disk media.
*/
func GetMediaPath(id string) (string, error) {
	// Grabbing media from database
	media, err := GetMedia(id)
	if err != nil {
		return "public/assets/default_pfp.jpg", err
	}

	return media.Path, nil
}
