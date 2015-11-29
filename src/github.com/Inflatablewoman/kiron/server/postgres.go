package server

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	// Required
	_ "github.com/lib/pq"
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

func (r postgresRepository) GetApplications() ([]*Application, error) {
	log.Printf("Going to get all applications. This is not gonna happen.")
	return nil, nil
}

func (r postgresRepository) GetApplicationsByStatus(status string) ([]*Application, error) {
	log.Printf("Going to get all applications for status %s", status)
	stmt, err := r.db.Prepare("SELECT id, birthday, phone, nationality, country, city, zip, address, address_extra, first_page_of_survey_data, gender, study_program, user_id, education_level_id, status, blocked_until, created_at, edited_at FROM applications WHERE status=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var apps []*Application
	for rows.Next() {
		app := &Application{}

		err = rows.Scan(&app.ID,
			&app.Birthday,
			&app.PhoneNumber,
			&app.Nationality,
			&app.Country,
			&app.City,
			&app.Zip,
			&app.Address,
			&app.AddressExtra,
			&app.FirstPageOfSurveyData,
			&app.Gender,
			&app.StudyProgram,
			&app.UserID,
			&app.EducationLevel,
			&app.Status,
			&app.BlockExpires,
			&app.Created,
			&app.Edited)
		if err != nil {
			log.Printf("Error with scan: %v", err)
			return nil, err
		}

		apps = append(apps, app)
	}

	return apps, nil
}

func (r postgresRepository) GetApplication(applicationID int) (*Application, error) {
	log.Printf("Going to application with id %d", applicationID)
	stmt, err := r.db.Prepare("SELECT id, birthday, phone, nationality, country, city, zip, address, address_extra, first_page_of_survey_data, gender, study_program, user_id, education_level_id, status, blocked_until, created_at, edited_at FROM applications WHERE id=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(applicationID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		app := &Application{}

		err = rows.Scan(&app.ID,
			&app.Birthday,
			&app.PhoneNumber,
			&app.Nationality,
			&app.Country,
			&app.City,
			&app.Zip,
			&app.Address,
			&app.AddressExtra,
			&app.FirstPageOfSurveyData,
			&app.Gender,
			&app.StudyProgram,
			&app.UserID,
			&app.EducationLevel,
			&app.Status,
			&app.BlockExpires,
			&app.Created,
			&app.Edited)
		if err != nil {
			log.Printf("Error with scan: %v", err)
			return nil, err
		}

		return app, nil
	}

	return nil, errors.New("Not found")
}

func (r postgresRepository) GetApplicationOf(userID int) (*Application, error) {
	log.Printf("Going to get application for user with id %d", userID)
	stmt, err := r.db.Prepare(`SELECT id, 
								birthday, 
								phone, 
								nationality, 
								country, 
								city, 
								zip,
								address, 
								address_extra, 
								first_page_of_survey_data, 
								gender, 
								study_program, 
								user_id, 
								education_level_id, 
								status, 
								blocked_until, 
								created_at, 
								edited_at 
								FROM applications WHERE user_id=$1`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		app := &Application{}

		err = rows.Scan(&app.ID,
			&app.Birthday,
			&app.PhoneNumber,
			&app.Nationality,
			&app.Country,
			&app.City,
			&app.Zip,
			&app.Address,
			&app.AddressExtra,
			&app.FirstPageOfSurveyData,
			&app.Gender,
			&app.StudyProgram,
			&app.UserID,
			&app.EducationLevel,
			&app.Status,
			&app.BlockExpires,
			&app.Created,
			&app.Edited)
		if err != nil {
			log.Printf("Error with scan: %v", err)
			return nil, err
		}

		return app, nil
	}

	return nil, errors.New("Not found")

}

func (r postgresRepository) SetApplication(application *Application) error {
	stmt, err := r.db.Prepare(`INSERT INTO applications 
								(birthday, 
								phone, 
								nationality, 
								country, 
								city, 
								zip, 
								address, 
								address_extra, 
								first_page_of_survey_data, 
								gender, 
								study_program, 
								user_id, 
								education_level_id, 
								status, 
								blocked_until, 
								created_at, 
								edited_at) 
								VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(
		application.Birthday,
		application.PhoneNumber,
		application.Nationality,
		application.Country,
		application.City,
		application.Zip,
		application.Address,
		application.AddressExtra,
		application.FirstPageOfSurveyData,
		application.Gender,
		application.StudyProgram,
		application.UserID,
		application.EducationLevel,
		application.Status,
		application.BlockExpires,
		application.Created,
		application.Edited)
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

func (r postgresRepository) UpdateApplication(application *Application) error {
	stmt, err := r.db.Prepare("UPDATE applications SET birthday=$1, phone=$2, nationality=$3, country=$4, city=$5, zip=$6, address=$7, address_extra=$8, first_page_of_survey_data=$9, gender=$10, user_id=$11, education_level_id=$12, status=$13, blocked_until=$14, created_at=$15, edited_at=$16 WHERE id=$17")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(
		application.Birthday,
		application.PhoneNumber,
		application.Nationality,
		application.Country,
		application.City,
		application.Zip,
		application.Address,
		application.AddressExtra,
		application.FirstPageOfSurveyData,
		application.Gender,
		application.StudyProgram,
		application.UserID,
		application.EducationLevel,
		application.Status,
		application.BlockExpires,
		application.Created,
		application.Edited,
		application.ID)
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

// we actually don't delete an application. Still, we need this function for Data Protection Law
func (r postgresRepository) DeleteApplication(applicationID int) error {
	stmt, err := r.db.Prepare("DELETE FROM applications WHERE id=$1")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(applicationID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Rows affected = %d\n", rowCnt)

	return nil
}

func (r postgresRepository) GetComments(applicationID int) ([]*Comment, error) {
	log.Printf("Going to get all comments for application id %d", applicationID)
	stmt, err := r.db.Prepare("SELECT id, created_at, user_id, contents FROM comments WHERE application_id=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(applicationID)
	if err != nil {
		return nil, err
	}

	var (
		commentID int
		createdAt time.Time
		userID    int
		contents  string
		comments  []*Comment
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&commentID, &createdAt, &userID, &contents)
		if err != nil {
			return nil, err
		}
		log.Println(commentID, createdAt, userID, contents)
		comment := Comment{ID: commentID, Created: createdAt, ApplicationID: applicationID, UserID: userID, Contents: contents}

		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r postgresRepository) GetComment(commentID int) (*Comment, error) {
	log.Printf("Going to get single comment with id %d", commentID)
	stmt, err := r.db.Prepare("SELECT created_at, application_id, user_id, contents FROM comments WHERE id=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(commentID)
	if err != nil {
		return nil, err
	}

	var (
		createdAt     time.Time
		applicationID int
		userID        int
		contents      string
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&createdAt, &applicationID, &userID, &contents)
		if err != nil {
			return nil, err
		}
		log.Println(createdAt, applicationID, userID, contents)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	comment := Comment{ID: commentID, Created: createdAt, ApplicationID: applicationID, UserID: userID, Contents: contents}

	return &comment, nil
}

func (r postgresRepository) SetComment(comment *Comment) error {
	stmt, err := r.db.Prepare("INSERT INTO comments(created_at, application_id, user_id, contents) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(comment.Created, comment.ApplicationID, comment.UserID, comment.Contents)
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

func (r postgresRepository) UpdateComment(comment *Comment) error {
	stmt, err := r.db.Prepare("UPDATE comments SET created_at=$1, application_id=$2, user_id=$3, contents=$4 WHERE id=$5")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(comment.Created, comment.ApplicationID, comment.UserID, comment.Contents, comment.ID)
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

func (r postgresRepository) DeleteComment(commentID int) error {
	stmt, err := r.db.Prepare("DELETE FROM comments WHERE id=$1")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(commentID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Rows affected = %d\n", rowCnt)

	return nil
}

func (r postgresRepository) GetUsers() ([]User, error) {
	return nil, nil
}

func (r postgresRepository) GetUser(userID int) (*User, error) {
	log.Printf("Going to get user by id: %v", userID)
	stmt, err := r.db.Prepare("SELECT id, email, name, lastname, password, created_at, role_id FROM users WHERE id=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var (
		id        int
		name      string
		email     string
		lastName  string
		password  string
		created   time.Time
		roleValue int
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &email, &name, &lastName, &password, &created, &roleValue)
		if err != nil {
			return nil, err
		}
		log.Println(id, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	user := User{ID: id, EmailAddress: email, FirstName: name, LastName: lastName, Password: password, Created: created, Role: role(roleValue)}

	return &user, nil
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

func (r postgresRepository) UpdateUser(user *User) error {

	stmt, err := r.db.Prepare("UPDATE users SET email=$1, name=$2, lastname=$3, password=$4, created_at=$5, role_id=$6 WHERE id=$7")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(user.EmailAddress, user.FirstName, user.LastName, user.Password, user.Created, user.Role.Value(), user.ID)
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

	log.Printf("Rows affected = %d\n", rowCnt)

	return nil
}

// Documents ...
func (r postgresRepository) GetDocuments(applicationID int) ([][]byte, error) {
	return nil, nil
}

func (r postgresRepository) StoreDocument(document *Document) error {
	stmt, err := r.db.Prepare("INSERT INTO documents(application_id, document_type_id, contents) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(document.ApplicationID, document.DocumentTypeID, document.Contents)
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

func (r postgresRepository) GetDocument(documentID int) (*Document, error) {
	log.Printf("Going to get document by ID:  %v", documentID)
	stmt, err := r.db.Prepare("SELECT application_id, document_type_id, contents FROM documents WHERE id=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(documentID)
	if err != nil {
		return nil, err
	}

	var (
		applicationID int
		docTypeID     int
		userID        int
		contents      []byte
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&applicationID, &docTypeID, &userID, &contents)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	document := Document{ID: documentID, ApplicationID: applicationID, DocumentTypeID: docTypeID, Contents: contents}

	return &document, nil
}

func (r postgresRepository) DeleteDocument(documentID int) error {
	log.Printf("Going to delete document with id = %d", documentID)
	stmt, err := r.db.Prepare("DELETE FROM documents WHERE id = $1")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(documentID)
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

func (r postgresRepository) GetToken(tokenValue string) (*Token, error) {
	log.Printf("Going to get token by value: %v", tokenValue)
	stmt, err := r.db.Prepare("SELECT user_id, expires FROM auth_tokens WHERE token=$1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(tokenValue)
	if err != nil {
		return nil, err
	}

	var (
		userID  int
		expires time.Time
	)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userID, &expires)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if userID == 0 {
		return nil, errors.New("Not found")
	}

	token := Token{UserID: userID, Value: tokenValue, Expires: expires}

	return &token, nil
}

func (r postgresRepository) SetToken(token *Token) error {
	stmt, err := r.db.Prepare("INSERT INTO auth_tokens(user_id, token, expires) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(token.UserID, token.Value, token.Expires)
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

func (r postgresRepository) DelToken(tokenValue string) error {
	log.Printf("Deleting token: %s", tokenValue)

	stmt, err := r.db.Prepare("DELETE FROM auth_tokens WHERE token = $1")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(tokenValue)
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

func (r postgresRepository) DelExpiredTokens() error {
	return nil
}
