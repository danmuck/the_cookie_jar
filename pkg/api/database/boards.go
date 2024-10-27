package database

import (
	"context"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a board to the database.
*/
func AddBoard(classroomID string, name string) error {
	// Grab classroom to add board to, also making sure classroom exists
	classroom, err := GetClassroom(classroomID)
	if err != nil {
		return err
	}

	// Creating the new board
	board := &models.Board{
		ID:        uuid.New().String(),
		Name:      name,
		ThreadIDs: make([]string, 0),
	}
	classroom.BoardIDs = append(classroom.BoardIDs, board.ID)

	// Trying to add board to database
	boardCollection := GetCollection("boards")
	_, err = boardCollection.InsertOne(context.TODO(), board)
	if err != nil {
		return err
	}

	return nil
}

/*
Gets a board model from the database.
*/
func GetBoard(id string) (*models.Board, error) {
	var board *models.Board
	err := GetCollection("boards").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&board)
	return board, err
}

/*
Will search for the clasroom in the database based on ID of given board
model and then update it to the given model.
*/
func UpdateBoard(board *models.Board) error {
	// Does the Board exist
	_, err := GetBoard(board.ID)
	if err != nil {
		return err
	}

	err = GetCollection("boards").FindOneAndReplace(context.TODO(), gin.H{"_id": board.ID}, board).Err()
	return err
}
