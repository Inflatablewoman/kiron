package server

import (
	"database/sql"
	_ "github.com/lib/pq"
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

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to postgres %v", err)
	}

	pr := postgresRepository{db: db}

	return pr, nil
}

func (r postgresRepository) GetApplications() ([]Application, error)                { return nil, nil }
func (r postgresRepository) GetApplication(applicationID int) (*Application, error) { return nil, nil }
func (r postgresRepository) SetApplication(application *Application) error          { return nil }
func (r postgresRepository) DeleteApplication(applicationID int) error              { return nil }

func (r postgresRepository) GetComments() ([]Comment, error)   { return nil, nil }
func (r postgresRepository) GetComment(commentID int) *Comment { return nil }
func (r postgresRepository) DeleteComment(commentID int) error { return nil }

func (r postgresRepository) GetUsers() ([]*User, error)           { return nil, nil }
func (r postgresRepository) GetUser(userID string) (*User, error) { return nil, nil }
func (r postgresRepository) DeleteUser(userID string) error       { return nil }
func (r postgresRepository) SetUser(*User) error                  { return nil }

func (r postgresRepository) StoreDocument(documentID string, data []byte) error { return nil }
func (r postgresRepository) GetDocument(documentID string) ([]byte, error)      { return nil, nil }
func (r postgresRepository) DeleteDocument(documentID string) error             { return nil }
