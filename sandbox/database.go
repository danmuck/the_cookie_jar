package sandbox

type Database struct {
	uri string
}

func NewDatabase() *Database {
	db := &Database{
		uri: "mongodb://localhost:27017",
	}
	return db
}
