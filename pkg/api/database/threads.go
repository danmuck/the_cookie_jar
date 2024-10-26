package database

import (
	"context"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a thread to the database.
*/
func AddThread(boardID string, name string) error {
	// Grab board to add thread to, also making sure board exists
	board, err := GetBoard(boardID)
	if err != nil {
		return err
	}

	// Creating the new thread
	thread := &models.Thread{
		ID:         uuid.New().String(),
		Name:       name,
		CommentIDs: make([]string, 0),
	}
	board.ThreadIDs = append(board.ThreadIDs, thread.ID)

	// Trying to add thread to database
	threadCollection := GetCollection("threads")
	_, err = threadCollection.InsertOne(context.TODO(), thread)
	if err != nil {
		return err
	}

	return nil
}

/*
Gets a thread model from the database.
*/
func GetThread(id string) (*models.Thread, error) {
	var thread *models.Thread
	err := GetCollection("threads").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&thread)
	return thread, err
}

/*
Will search for the clasroom in the database based on ID of given thread
model and then update it to the given model.
*/
func UpdateThread(thread *models.Thread) error {
	// Does the Thread exist
	_, err := GetThread(thread.ID)
	if err != nil {
		return err
	}

	err = GetCollection("threads").FindOneAndReplace(context.TODO(), gin.H{"_id": thread.ID}, thread).Err()
	return err
}
