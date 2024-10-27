package database

import (
	"context"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a classroom to the database.
*/
func AddClassroom(username string, name string) error {
	// Grab user who is making classroom, also making sure they exist
	user, err := GetUser(username)
	if err != nil {
		return err
	}

	// Creating the new classroom
	classroom := &models.Classroom{
		ID:            uuid.New().String(),
		Name:          name,
		ProfessorID:   user.ID,
		InstructorIDs: make([]string, 0),
		StudentIDs:    make([]string, 0),
		BoardIDs:      make([]string, 0),
	}

	// Trying to add classroom to database
	classroomCollection := GetCollection("classrooms")
	_, err = classroomCollection.InsertOne(context.TODO(), classroom)
	if err != nil {
		return err
	}

	return nil
}

/*
Gets a classroom model from the database.
*/
func GetClassroom(id string) (*models.Classroom, error) {
	var classroom *models.Classroom
	err := GetCollection("classrooms").FindOne(context.TODO(), gin.H{"_id": id}).Decode(&classroom)
	return classroom, err
}

/*
Will search for the clasroom in the database based on ID of given classroom
model and then update it to the given model.
*/
func UpdateClassroom(classroom *models.Classroom) error {
	// Does the classroom exist
	_, err := GetClassroom(classroom.ID)
	if err != nil {
		return err
	}

	err = GetCollection("classrooms").FindOneAndReplace(context.TODO(), gin.H{"_id": classroom.ID}, classroom).Err()
	return err
}
