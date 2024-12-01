package database

import (
	"context"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a thread to the database.
*/
func AddThread(username string, title string, classroomId string, comment string) (*models.Thread, error) {
	// Making sure the user exists
	_, err := GetUser(username)
	if err != nil {
		return nil, err
	}

	// Creating the new thread
	thread := &models.Thread{
		ID:         uuid.New().String(),
		Title:      title,
		Date:       time.Now().Format("01/02/2006"),
		AuthorID:   username,
		CommentIDs: make([]string, 0),
	}

	// Trying to add thread to database
	threadCollection := GetCollection("threads")
	_, err = threadCollection.InsertOne(context.TODO(), thread)
	if err != nil {
		return nil, err
	}

	// Associating thread with the classroom
	err = UpdateClassroomThreads(classroomId, thread.ID)
	if err != nil {
		return nil, err
	}

	// Adding the first comment to the thread
	_, err = AddComment(username, comment, thread.ID)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

/*
Will search for the thread in the database and then add a comment in it.
*/
func UpdateThreadComments(threadId string, commentId string) error {
	// Grab the thread and verify it exists
	thread, err := GetThread(threadId)
	if err != nil {
		return err
	}

	thread.CommentIDs = append(thread.CommentIDs, commentId)
	err = GetCollection("threads").FindOneAndReplace(context.TODO(), gin.H{"_id": threadId}, thread).Err()
	return err
}

/*
Gets a thread model from the database.
*/
func GetThread(id string) (*models.Thread, error) {
	var thread *models.Thread
	err := GetCollection("threads").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&thread)
	return thread, err
}
