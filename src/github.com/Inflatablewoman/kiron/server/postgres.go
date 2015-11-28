package server

import (
	"database/sql"
	"os"
	"time"
	// Required
	_ "github.com/lib/pq"
	"log"
)

type postgresRepository struct {
	db *sql.DB
}

func getPostgresDB() (DataRepository, error) {

	connectionString := os.Getenv("DBCONN")

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
	stmt, err := r.db.Prepare("SELECT id, name, lastname, password, created_at, role_id FROM users WHERE email=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(emailAddress)
	if err != nil {
		return nil, err
	}

	var (
		id        int
		name      string
		lastName  string
		password  string
		created   time.Time
		roleValue int
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &lastName, &password, &created, &roleValue)
		if err != nil {
			return nil, err
		}
		log.Println(id, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	user := User{ID: id, EmailAddress: emailAddress, FirstName: name, LastName: lastName, Password: password, Created: created, Role: role(roleValue)}

	return &user, nil
}

func (r postgresRepository) DeleteUser(userID int) error {
	stmt, err := r.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(userID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("Got error - LastInsertId: %v", err)
	}
	log.Printf("affected = %d\n", rowCnt)

	return nil
}

func (r postgresRepository) SetUser(user *User) error {

	stmt, err := r.db.Prepare("INSERT INTO users(email, name, lastname, password, created_at, role_id) VALUES($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(user.EmailAddress, user.FirstName, user.LastName, user.Password, user.Created, user.Role.Value())
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("Got error - RowsAffected: %v", err)
	}
	log.Printf("affected = %d\n", rowCnt)

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
