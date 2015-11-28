package server

var repository DataRepository

// InitDatabase will create a connection to a postgres server
func InitDatabase(connectionString string) error {
	var err error
	repository, err = getPostgresDB(connectionString)
	return err
}

// DataRepository is a repository
type DataRepository interface {
	GetUser(userID string) (*User, error)
	DeleteUser(userID string) error
	SetUser(*User) error
	StoreDocument(documentID string, data []byte) error
	GetDocument(documentID string) ([]byte, error)
}
