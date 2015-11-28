package server

import (
	"database/sql"
	"log"
)

type postgresRepository struct {
	db *sql.DB
}

func getPostgresDB(connectionString string) (DataRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to postgres %v", err)
	}
	pr := postgresRepository{db: db}
	return pr, nil
}

func (r postgresRepository) GetUser(userID string) (*User, error) {
	return &User{}, nil
}

func (r postgresRepository) DeleteUser(userID string) error {
	return nil
}

func (r postgresRepository) SetUser(*User) error {
	return nil
}

func (r postgresRepository) StoreDocument(documentID string, data []byte) error {
	return nil
}

func (r postgresRepository) GetDocument(documentID string) ([]byte, error) {
	return []byte{}, nil
}
