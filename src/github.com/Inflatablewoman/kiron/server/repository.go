package server

import (
	"log"
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
	UpdateApplication(application *Application) error
	DeleteApplication(applicationID int) error

	GetComments(applicationID int) ([]*Comment, error)
	GetComment(commentID int) (*Comment, error)
	SetComment(comment *Comment) error
	UpdateComment(comment *Comment) error
	DeleteComment(commentID int) error

	GetUsers() ([]User, error)
	GetUser(userID int) (*User, error)
	GetUserByEmail(emailAddress string) (*User, error)
	SetUser(*User) error
	UpdateUser(*User) error
	DeleteUser(userID int) error

	GetDocuments(userID string, applicationID string) ([][]byte, error)
	StoreDocument(documentID string, data []byte) error
	GetDocument(documentID string) ([]byte, error)
	DeleteDocument(documentID string) error

	GetToken(tokenValue string) (*Token, error)
	SetToken(token *Token) error
	DelToken(tokenValue string) error
	DelExpiredTokens()
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

// ToRestUser converts repo version of User to RestUser
func (u *User) ToRestUser() *RestUser {
	ru := RestUser{ID: u.ID, EmailAddress: u.EmailAddress, FirstName: u.FirstName, LastName: u.LastName, Created: u.Created, Role: u.Role}
	return &ru
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

// ToRestApplication converts repo version of Application to RestApplication
func (a *Application) ToRestApplication() *RestApplication {
	user, err := repository.GetUser(a.UserID)

	if err != nil {
		log.Println("An application without user reference found.")
		return nil
	}

	ru := RestApplication{
		ID: a.ID, UserID: a.UserID, FirstName: user.FirstName,
		LastName: user.LastName, Birthday: a.Birthday, PhoneNumber: a.PhoneNumber,
		Nationality: a.Nationality, Address: a.Address, AddressExtra: a.AddressExtra,
		Zip: a.Zip, City: a.City, Country: a.Country, FirstPageOfSurveyData: a.FirstPageOfSurveyData,
		Gender: a.Gender, EducationLevel: a.EducationLevel, Status: a.Status,
		Created: a.Created, Edited: a.Edited}
	return &ru
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
	Token       string
	TokenExpiry int
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

// Token ...
type Token struct {
	UserID  int
	Value   string
	Expires time.Time
}
