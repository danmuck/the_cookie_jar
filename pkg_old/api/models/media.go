package models

import "time"

type Media struct {
    ID        string    `bson:"_id" json:"id"`
    Username  string    `bson:"username" json:"username"`
    Format    string    `bson:"format" json:"format"`
    Path      string    `bson:"path" json:"path"`
    Size      int64     `bson:"size" json:"size"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
