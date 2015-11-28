package server

import (
	"database/sql"
	// Required
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

func (r postgresRepository) GetApplications() ([]Application, error) {
	return nil, nil
}

func (r postgresRepository) GetApplication(applicationID int) (*Application, error) {
	return nil, nil
}
func (r postgresRepository) SetApplication(application *Application) error {
	return nil
}
func (r postgresRepository) DeleteApplication(applicationID int) error {
	return nil
}

func (r postgresRepository) GetComments() ([]Comment, error) {
	return nil, nil
}

func (r postgresRepository) GetComment(commentID int) *Comment {
	return nil
}

func (r postgresRepository) DeleteComment(commentID int) error {
	return nil
}

func (r postgresRepository) GetUses() ([]*User, error) {
	return nil, nil
}

func (r postgresRepository) GetUser(userID string) (*User, error) {
	return nil, nil
}

func (r postgresRepository) GetUserByEmail(emailAddress string) (*User, error) {
	log.Printf("Going to get user by email address %v", emailAddress)
	stmt, err := r.db.Prepare("SELECT id, name FROM users WHERE email=$N")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(emailAddress)
	if err != nil {
		return nil, err
	}

	var (
		id   int
		name string
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	user := User{ID: id, EmailAddress: emailAddress, FirstName: name}

	return &user, nil
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
	return nil, nil
}
func (r postgresRepository) DeleteDocument(documentID string) error {
	return nil
}
