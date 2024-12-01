package database

import (
	"context"
	"fmt"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a comment to the database.
*/
func AddComment(username string, text string, threadId string) (*models.Comment, error) {
	// Making sure the user exists
	_, err := GetUser(username)
	if err != nil {
		return nil, err
	}

	// Creating the new comment
	comment := &models.Comment{
		ID:           uuid.New().String(),
		Text:         text,
		AuthorID:     username,
		LikedUserIDs: make([]string, 0),
	}

	// Trying to add comment to database
	commentCollection := GetCollection("comments")
	_, err = commentCollection.InsertOne(context.TODO(), comment)
	if err != nil {
		return nil, err
	}

	// Associating comment with the thread
	err = UpdateThreadComments(threadId, comment.ID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

/*
Returns true/false if a user liked a comment ID, also does the opposite action
of what they did.
*/
func IsUserLiked(id string, username string) (bool, error) {
	comment, err := GetComment(id)
	if err != nil {
		return false, err
	}

	isLiked := utils.Contains(comment.LikedUserIDs, username)
	if isLiked {
		comment.LikedUserIDs = utils.RemoveItem(comment.LikedUserIDs, username)
	} else {
		comment.LikedUserIDs = append(comment.LikedUserIDs, username)
	}
	err = GetCollection("comments").FindOneAndReplace(context.TODO(), gin.H{"_id": id}, comment).Err()

	return !isLiked, err
}

/*
Returns nil/error if the given username is the author of the given comment ID.
*/
func IsUsersComment(id string, username string) error {
	comment, err := GetComment(id)
	if err != nil {
		return err
	}

	if comment.AuthorID == username {
		err = nil
	} else {
		err = fmt.Errorf("not this users comment")
	}
	return err
}

/*
Will search for the comment in the database and then change its text.
*/
func UpdateCommentText(id string, text string) error {
	// Grab the comment and verify it exists
	comment, err := GetComment(id)
	if err != nil {
		return err
	}

	comment.Text = text
	err = GetCollection("comments").FindOneAndReplace(context.TODO(), gin.H{"_id": id}, comment).Err()
	return err
}

/*
Gets a comment model from the database.
*/
func GetComment(id string) (*models.Comment, error) {
	var comment *models.Comment
	err := GetCollection("comments").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&comment)
	return comment, err
}
