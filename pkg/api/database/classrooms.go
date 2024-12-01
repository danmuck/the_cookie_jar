package database

import (
	"context"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/danmuck/the_cookie_jar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Adds a classroom to the database.
*/
func AddClassroom(username string, name string) error {
	// Making sure the user exists
	_, err := GetUser(username)
	if err != nil {
		return err
	}

	// Creating the new classroom
	classroom := &models.Classroom{
		ID:          uuid.New().String(),
		Name:        name,
		ProfessorID: username,
		StudentIDs:  make([]string, 0),
		ThreadIDs:   make([]string, 0),
		Game:        models.ClassGame{},
	}

	// Trying to add classroom to database
	classroomCollection := GetCollection("classrooms")
	_, err = classroomCollection.InsertOne(context.TODO(), classroom)
	if err != nil {
		return err
	}

	// Associating user with the classroom directly
	err = UpdateUserJoinedClassrooms(username, classroom.ID, false)
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
Will search for the classroom and user in the database and then add/remove a
classroom for them.
*/
func UpdateClassroomStudents(classroomId string, username string, remove bool) error {
	// Grab the classroom and verify it exists
	classroom, err := GetClassroom(classroomId)
	if err != nil {
		return err
	}

	// Adding/removing classroom from the user
	err = UpdateUserJoinedClassrooms(username, classroomId, remove)
	if err != nil {
		return err
	}

	if remove {
		classroom.StudentIDs = utils.RemoveItem(classroom.StudentIDs, username)
	} else {
		classroom.StudentIDs = append(classroom.StudentIDs, username)
	}
	err = GetCollection("classrooms").FindOneAndReplace(context.TODO(), gin.H{"_id": classroomId}, classroom).Err()
	return err
}

/*
Will search for the classroom in the database and then add a thread in it.
*/
func UpdateClassroomThreads(classroomId string, threadId string) error {
	// Grab the classroom and verify it exists
	classroom, err := GetClassroom(classroomId)
	if err != nil {
		return err
	}

	classroom.ThreadIDs = append(classroom.ThreadIDs, threadId)
	err = GetCollection("classrooms").FindOneAndReplace(context.TODO(), gin.H{"_id": classroomId}, classroom).Err()
	return err
}
