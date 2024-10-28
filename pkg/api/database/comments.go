package database

import (
	"context"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a comment to the database.
*/
func AddComment(threadID string, username string, text string) error {
	// Grab thread to add comment to, also making sure thread exists
	thread, err := GetThread(threadID)
	if err != nil {
		return err
	}

	// Grab username of user who commented, also making sure user exists
	user, err := GetUser(username)
	if err != nil {
		return err
	}

	// Creating the new comment
	comment := &models.Comment{
		ID:           uuid.New().String(),
		Username:     user.Username,
		Text:         text,
		LikedUserIDs: make([]string, 0),
		Date:         time.Now(),
	}

	// Trying to add comment to database
	commentCollection := GetCollection("comments")
	_, err = commentCollection.InsertOne(context.TODO(), comment)
	if err != nil {
		return err
	}

	// Associating thread with the comment directly
	thread.CommentIDs = append(thread.CommentIDs, comment.ID)
	err = UpdateThread(thread)
	if err != nil {
		return err
	}

	

	return nil
}

/*
Will create stats for a user if they don't exist.
*/
func CreateUserStats(username string) error {
	stats := &models.UserCommentStats{
		ID:				uuid.New().String(),
		Username:		username,
		TotalComments:	0,
		TotalLikes:		0,
	}
	_, err := GetCollection("user_stats").InsertOne(context.TODO(), stats)
	return err
}

/*
This will retireve the stats for a given 
*/
func GetUserStats(username string) (*models.UserCommentStats, error) {
    var stats models.UserCommentStats
    err := GetCollection("user_stats").FindOne(context.TODO(), gin.H{"username": username},).Decode(&stats)
    return &stats, err
}
/*
Will update the user stats for an associated user.
*/

func UpdateUserStats(username string) error {
    _, err := GetCollection("user_stats").UpdateOne(
        context.TODO(),
        gin.H{"username": username},
        gin.H{
            "$inc": gin.H{"total_comments": 1},
        },
    )
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

/*
Will search for the clasroom in the database based on ID of given comment
model and then update it to the given model.
*/
func UpdateComment(comment *models.Comment) error {
	// Does the Comment exist
	_, err := GetComment(comment.ID)
	if err != nil {
		return err
	}

	err = GetCollection("comments").FindOneAndReplace(context.TODO(), gin.H{"_id": comment.ID}, comment).Err()
	return err
}
