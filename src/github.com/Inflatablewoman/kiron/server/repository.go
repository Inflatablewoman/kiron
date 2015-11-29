package server

import (
	"time"
)

var repository DataRepository

// InitDatabase will create a connection to a postgres server
func InitDatabase() error {
	var err error
	repository, err = getPostgresDB()
	return err
}

// DataRepository is a repository
type DataRepository interface {
	GetApplications() ([]*Application, error) // this cannot return all applications, at least not with files.
	GetApplicationsByStatus(status string) ([]*Application, error)
	GetApplication(applicationID int) (*Application, error)
	GetApplicationOf(userID int) (*Application, error)
	SetApplication(application *Application) error
	DeleteApplication(applicationID int) error

	GetComments(applicationID int) ([]*Comment, error)
	SetComment(comment *Comment) error
	GetComment(commentID int) (*Comment, error)
	DeleteComment(commentID int) error

	GetUsers() ([]User, error)
	GetUser(userID string) (*User, error)
	GetUserByEmail(emailAddress string) (*User, error)
	DeleteUser(userID int) error
	SetUser(*User) error

	GetDocuments(userID string, applicationID string) ([][]byte, error)
	StoreDocument(documentID string, data []byte) error
	GetDocument(documentID string) ([]byte, error)
	DeleteDocument(documentID string) error
}

// Roles ...
type role int

// RoleNone Define Constant stati
const RoleNone role = 0

// Other roles
const (
	RoleAdmin role = 1 << iota
	RoleSubAdmin
	RoleTrustedHelper
	RoleLimitedHelper
	RoleApplication
)

// Value return the value of the role
func (r role) Value() int {
	return int(r)
}

// User is the users struct
type User struct {
	ID           int
	EmailAddress string
	FirstName    string
	LastName     string
	Password     string
	Created      time.Time
	Role         role
}

// Status ...
type status string

const statusNone = 0

const (
	statusRecieved role = 1 << iota
	statusConfirmed
	statusInVerification
	statusVerified
	statusWaitingForResponse
	statusRejected
	statusAccepted
)

// ALlowed stati
var allowedStati = []string{"received", "confirmed", "inVerfication", "verified", "waitingForResponse", "rejected", "accepted"}

// Allowed education types
var allowedEducationTypes = []string{"none", "elementary", "secondary", "associate", "bachelor"}

// Application ...
type Application struct {
	ID                    int
	Birthday              time.Time
	PhoneNumber           string
	Nationality           string
	Country               string
	City                  string
	Zip                   string
	Address               string
	AddressExtra          string
	FirstPageOfSurveyData string
	Gender                string
	UserID                int
	EducationLevel        int
	Status                string
	BlockExpires          time.Time
	Created               time.Time
	Edited                time.Time
}

// Comment ...
type Comment struct {
	ID            int
	Created       time.Time
	ApplicationID int
	UserID        int
	Contents      string
}

// Document ...
type Document struct {
	ID             int
	ApplicationID  int
	DocumentTypeID int
	UserID         int
	Contents       string
}

// LoginResponse ...
type LoginResponse struct {
	Token       int
	TokenExpory int
	Result      LoginResult
}

// LoginResult  ...
type LoginResult struct {
	ID           int
	EmailAddress string
	FirstName    string
	LastName     string
	Role         role
}
