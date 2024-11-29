package models

type Media struct {
	ID               string `bson:"_id" json:"id"`
	UploaderUsername string `bson:"username" json:"username"`
	Format           string `bson:"format" json:"format"`
}
