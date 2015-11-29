package server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rcrowley/go-tigertonic"
)

// AuthContext information for Marshaled calls
type AuthContext struct {
	UserAgent  string
	RemoteAddr string
	User       *User
	TokenValue string
}

// getContext is used check the Auth of a user
func getContext(r *http.Request) (http.Header, error) {

	tigertonic.Context(r).(*AuthContext).UserAgent = r.UserAgent()
	tigertonic.Context(r).(*AuthContext).RemoteAddr = RequestAddr(r)

	tokenValue := GetBearerAuthFromHeader(r.Header)

	tigertonic.Context(r).(*AuthContext).TokenValue = tokenValue

	if tokenValue == "" {
		log.Println("No token passed")
		return nil, tigertonic.Unauthorized{Err: errors.New("Access denied")}
	}

	// Check token is valid
	token, err := repository.GetToken(tokenValue)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		return nil, tigertonic.Unauthorized{Err: errors.New("Access denied")}
	}

	if token == nil || token.Value == "" {
		log.Println("Token not found")
		return nil, tigertonic.Unauthorized{Err: errors.New("Access denied")}
	}

	// Get user
	user, err := repository.GetUser(token.UserID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, tigertonic.Unauthorized{Err: errors.New("Access denied")}
	}

	if user.ID == 0 {
		log.Println("User not found")
		return nil, tigertonic.Unauthorized{Err: errors.New("Access denied")}
	}

	// All good.  Add user to context
	tigertonic.Context(r).(*AuthContext).User = user

	return nil, nil
}

// BasicContext information for Marshaled calls
type BasicContext struct {
	UserAgent  string
	RemoteAddr string
}

// getContext is used to return UserAgent and Request info from the request
func getBasicContext(r *http.Request) (http.Header, error) {

	tigertonic.Context(r).(*BasicContext).UserAgent = r.UserAgent()
	tigertonic.Context(r).(*BasicContext).RemoteAddr = RequestAddr(r)

	return nil, nil
}

// RegisterHTTPHandlers registers the http handlers
func RegisterHTTPHandlers(mux *tigertonic.TrieServeMux) {

	// Login User
	mux.Handle("POST", "/api/v1/login", tigertonic.WithContext(tigertonic.If(getBasicContext, tigertonic.Marshaled(login)), BasicContext{}))

	// Logout
	mux.Handle("POST", "/api/v1/logout", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(logout)), AuthContext{}))

	// Create user
	mux.Handle("POST", "/api/v1/users", tigertonic.WithContext(tigertonic.If(getBasicContext, tigertonic.Marshaled(createUser)), BasicContext{}))

	// Get users
	mux.Handle("GET", "/api/v1/users", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getUsers)), AuthContext{}))

	// Get single user
	mux.Handle("GET", "/api/v1/users/{userID}", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getUser)), AuthContext{}))

	// Get applications
	mux.Handle("GET", "/api/v1/applications", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getApplications)), AuthContext{}))

	// Get single application
	mux.Handle("GET", "/api/v1/users/{userID}/application", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getApplication)), AuthContext{}))

	// Create application
	mux.Handle("POST", "/api/v1/users/{userID}/application", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(createApplication)), AuthContext{}))

	// Get documents
	mux.Handle("GET", "/api/v1/users/{userID}/application/{applicationID}/documents", NewFileDownloadHandler())

	// Create documents
	mux.Handle("PUT", "/api/v1/users/{userID}/application/{applicationID}/documents", NewRawUploadHandler())

	// Get comments
	mux.Handle("GET", "/api/v1/users/{userID}/application/{applicationID}/comments", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getComments)), AuthContext{}))

	// Create comment
	mux.Handle("POST", "/api/v1/users/{userID}/application/{applicationID}/comments", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(createComment)), AuthContext{}))

}

type loginRequest struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

// login
func login(u *url.URL, h http.Header, request *loginRequest, context *BasicContext) (int, http.Header, *LoginResponse, error) {
	var err error
	defer CatchPanic(&err, "login")

	log.Printf("login called: %s %s", context.RemoteAddr, context.UserAgent)

	if request.EmailAddress == "" {
		log.Println("Error:  No email address")
		return http.StatusBadRequest, nil, nil, errors.New("You must provide an email address")
	}

	if request.Password == "" {
		log.Println("Error:  No password")
		return http.StatusBadRequest, nil, nil, errors.New("You must provide a password")
	}

	user, err := repository.GetUserByEmail(request.EmailAddress)
	if err != nil {
		log.Printf("Error:  Unable to get user from repo: %v", err)
		return http.StatusUnauthorized, nil, nil, errors.New("Invalid password or unknown user")
	}

	match, _ := MatchPassword(request.Password, user.Password)
	if !match {
		log.Println("Error:  Password is incorrect")
		return http.StatusUnauthorized, nil, nil, errors.New("Invalid password or unknown user")
	}

	tokenValue := GetRandomString(16, "")
	expires := time.Now().UTC().Add(time.Duration(1 * time.Hour))

	t := Token{UserID: user.ID, Value: tokenValue, Expires: expires}

	// Save token to database
	err = repository.SetToken(&t)
	if err != nil {
		log.Printf("Error:  Unable to get token from repo: %v", err)
		return http.StatusInternalServerError, nil, nil, errors.New("Internal error")
	}

	expiresSeconds := 3600 // seconds in 1 hour
	lr := LoginResult{ID: user.ID,
		EmailAddress: user.EmailAddress,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role}
	lResp := LoginResponse{Token: tokenValue,
		TokenExpiry: expiresSeconds,
		Result:      lr}

	// All good!
	return http.StatusOK, nil, &lResp, nil
}

// logout will logout a session
func logout(u *url.URL, h http.Header, _ *emptyRequest, context *AuthContext) (int, http.Header, interface{}, error) {
	var err error
	defer CatchPanic(&err, "logout")

	log.Printf("Going to logout: %v", context.User.EmailAddress)

	err = repository.DelToken(context.TokenValue)
	if err != nil {
		log.Printf("Error deleting token: %v", err)
		return http.StatusInternalServerError, nil, nil, errors.New("Unable to delete token")
	}

	return http.StatusOK, nil, nil, nil
}

// RestUser ...
type RestUser struct {
	ID           int       `json:"id"`
	EmailAddress string    `json:"email"`
	FirstName    string    `json:"name"`
	LastName     string    `json:"lastname"`
	Created      time.Time `json:"created"`
	Role         role      `json:"role"`
}

type createUserRequest struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	LastName     string `json:"lastname"`
}

// createUser will create a user
func createUser(u *url.URL, h http.Header, request *createUserRequest, context *BasicContext) (int, http.Header, *RestUser, error) {
	var err error
	defer CatchPanic(&err, "createUser")

	log.Printf("createUser called by: %s %s", context.RemoteAddr, context.UserAgent)

	hashedPassword, _ := createHashedPassword(request.Password)
	user := User{EmailAddress: request.EmailAddress, Password: hashedPassword, FirstName: request.Name, LastName: request.LastName, Role: RoleAdmin}

	err = repository.SetUser(&user)

	if err != nil {
		log.Printf("Unable to set user: %v", err)
		return http.StatusInternalServerError, nil, nil, nil
	}

	ru := user.ToRestUser()

	// All good!
	return http.StatusOK, nil, ru, nil
}

// getUser will get a user
func getUser(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, *RestUser, error) {
	var err error
	defer CatchPanic(&err, "getUser")

	log.Println("getUser Started")

	userID, err := strconv.Atoi(u.Query().Get("userID"))

	// logged in applicant trying to view other user's profile
	if context.User.Role == RoleApplication && context.User.ID != userID {
		return http.StatusUnauthorized, nil, nil, errors.New("Access denied")
	}

	user, err := repository.GetUser(userID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, user.ToRestUser(), nil
}

// getUsers will get a user
func getUsers(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, []User, error) {
	var err error
	defer CatchPanic(&err, "getUsers")

	log.Println("getUsers Started")

	if context.User.Role != RoleAdmin && context.User.Role != RoleSubAdmin {
		return http.StatusUnauthorized, nil, nil, errors.New("Access denied")
	}

	users, err := repository.GetUsers()

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, users, nil
}

// getApplications will get a list of all (???) applications
func getApplications(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, []*RestApplication, error) {
	var err error
	defer CatchPanic(&err, "getApplications")

	log.Println("getApplications Started")

	// Logged in applicant is trying to access a list of applications
	if context.User.Role == RoleApplication {
		return http.StatusUnauthorized, nil, nil, errors.New("Access denied")
	}

	applications, err := repository.GetApplications()

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	var restApplications []*RestApplication

	for _, application := range applications {
		restApplication := application.ToRestApplication()

		if context.User.Role == RoleLimitedHelper {
			trimForLimitedHelper(restApplication)
		}

		restApplications = append(restApplications, restApplication)
	}

	// All good!
	return http.StatusOK, nil, restApplications, nil
}

// getApplication will get an application
func getApplication(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, *RestApplication, error) {
	var err error
	defer CatchPanic(&err, "getApplication")

	log.Println("getApplication Started")

	userID, err := strconv.Atoi(u.Query().Get("userID"))

	// Logged in applicant is trying to view another user's application
	if context.User.Role == RoleApplication && context.User.ID != userID {
		return http.StatusUnauthorized, nil, nil, errors.New("Access denied")
	}

	application, err := repository.GetApplicationOf(userID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	restApplication := application.ToRestApplication()

	if context.User.Role == RoleLimitedHelper {
		trimForLimitedHelper(restApplication)
	}

	// All good!
	return http.StatusOK, nil, restApplication, nil
}

// RestApplication ...
type RestApplication struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"user_id"`
	FirstName             string    `json:"name"`
	LastName              string    `json:"lastname"`
	Birthday              time.Time `json:"birthday"`
	PhoneNumber           string    `json:"phone"`
	Nationality           string    `json:"nationality"`
	Address               string    `json:"address"`
	AddressExtra          string    `json:"address_extra"`
	Zip                   string    `json:"zip"`
	City                  string    `json:"city"`
	Country               string    `json:"country"`
	FirstPageOfSurveyData string    `json:"first_page_of_survey_data"`
	Gender                string    `json:"gender"`
	EducationLevel        int       `json:"education_level_id"`
	Status                string    `json:"status"`
	Created               time.Time `json:"created_at"`
	Edited                time.Time `json:"edited_at"`
}

// trimForLimitedHelper trims some fields so that limited helper cannot view them
func trimForLimitedHelper(application *RestApplication) *RestApplication {
	application.PhoneNumber = ""
	application.Nationality = ""
	application.Address = ""
	application.AddressExtra = ""
	application.Zip = ""
	application.City = ""
	application.Country = ""
	application.FirstPageOfSurveyData = ""
	application.Gender = ""
	application.EducationLevel = 0
	return application
}

type createApplicationRequest struct {
	Birthday              time.Time `json:"birthday"`
	PhoneNumber           string    `json:"phone"`
	Nationality           string    `json:"nationality"`
	Address               string    `json:"address"`
	AddressExtra          string    `json:"address_extra"`
	Zip                   string    `json:"zip"`
	City                  string    `json:"city"`
	Country               string    `json:"country"`
	FirstPageOfSurveyData string    `json:"first_page_of_survey_data"`
	Gender                string    `json:"gender"`
	EducationLevel        int       `json:"education_level_id"`
}

func createApplication(u *url.URL, h http.Header, request *createApplicationRequest, context *AuthContext) (int, http.Header, *RestApplication, error) {
	var err error
	defer CatchPanic(&err, "createApplications")

	log.Println("createApplications Started")

	userID, err := strconv.Atoi(u.Query().Get("userID"))

	application := Application{Birthday: request.Birthday, PhoneNumber: request.PhoneNumber, Nationality: request.Nationality, Country: request.Country, City: request.City, Zip: request.Zip, AddressExtra: request.AddressExtra, FirstPageOfSurveyData: request.FirstPageOfSurveyData, Gender: request.Gender, UserID: userID, EducationLevel: request.EducationLevel}

	err = repository.SetApplication(&application)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusCreated, nil, application.ToRestApplication(), nil
}

func getDocuments(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, [][]byte, error) {
	var err error
	defer CatchPanic(&err, "getDocuments")

	log.Println("getDocuments Started")

	applicationID, err := strconv.Atoi(u.Query().Get("applicationID"))

	documents, err := repository.GetDocuments(applicationID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, documents, nil
}

func createDocuments(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, []*Document, error) {
	var err error
	defer CatchPanic(&err, "createDocuments")

	log.Println("createDocuments Started")

	// TODO Implement functionality
	//documents := nil

	// All good!
	return http.StatusOK, nil, nil, nil
}

func getComments(u *url.URL, h http.Header, _ interface{}, context *AuthContext) (int, http.Header, []*Comment, error) {
	var err error
	defer CatchPanic(&err, "getComments")

	log.Println("getComments Started")

	applicationID, err := strconv.Atoi(u.Query().Get("applicationID"))

	comments, err := repository.GetComments(applicationID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, comments, nil
}

type createCommentRequest struct {
	UserID   int    `json:"user_id"`
	Contents string `json:"contents"`
}

func createComment(u *url.URL, h http.Header, request *createCommentRequest, context *AuthContext) (int, http.Header, *Comment, error) {
	var err error
	defer CatchPanic(&err, "createComment")

	log.Println("createComment Started")

	applicationID, err := strconv.Atoi(u.Query().Get("applicationID"))
	comment := Comment{ApplicationID: applicationID, UserID: request.UserID, Contents: request.Contents}

	err = repository.SetComment(&comment)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, nil, nil
}

// RawUploadHandler handles PUT operations
type RawUploadHandler struct {
}

// NewRawUploadHandler ...
func NewRawUploadHandler() RawUploadHandler {
	return RawUploadHandler{}
}

// ServeHTTP ...
func (handler RawUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	defer CatchPanic(&err, "BlockRawUploadHandler")
	log.Println("Got PUT upload request")

	_, err = getContext(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Can I do this?
	//context := tigertonic.Context(r).(*AuthContext)

	applicationID, err := strconv.Atoi(r.Header.Get("applicationID"))
	documentTypeID, err := strconv.Atoi(r.Header.Get("documentTypeID"))

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleErrorWithResponse(w, err)
		return
	}

	document := Document{ApplicationID: applicationID, DocumentTypeID: documentTypeID, Contents: body}

	err = repository.StoreDocument(&document)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// FileDownloadHandler ...
type FileDownloadHandler struct {
}

// NewFileDownloadHandler ...
func NewFileDownloadHandler() FileDownloadHandler {
	return FileDownloadHandler{}
}

// ServeHTTP ...
func (handler FileDownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	defer CatchPanic(&err, "FileDownloadHandler")

	_, err = getContext(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	documentID, err := strconv.Atoi(r.URL.Query().Get("documentID"))
	if err != nil {
		HandleErrorWithResponse(w, err)
		return
	}

	document, err := repository.GetDocument(documentID)
	if err != nil {
		HandleErrorWithResponse(w, err)
		return
	}

	log.Printf("Download Document [%v]", documentID)

	header := w.Header()
	header["Content-Type"] = []string{"application/octet-stream"}
	w.Write(document.Contents)

	log.Println("Download complete")
}

// checkClose is used to check the return from Close in a defer
// statement.
func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}

// HandleErrorWithResponse ...
func HandleErrorWithResponse(w http.ResponseWriter, err error) {
	log.Printf("Got error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err)
	return
}

// CatchPanic will handle a panicm log an error and recover
func CatchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		log.Printf("%s : PANIC Defered : %v\n", functionName, r)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		log.Printf("%s : Stack Trace : %s", functionName, string(buf))
	}
}

// RequestAddr Get the request address
func RequestAddr(r *http.Request) string {
	// Get the IP of the request
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// GetBearerAuthFromHeader returns "Bearer" token from request.
func GetBearerAuthFromHeader(h http.Header) string {
	authHeader := h.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	s := strings.SplitN(authHeader, " ", 2)
	if len(s) != 2 || s[0] != "Bearer" {
		return ""
	}
	return s[1]
}

type emptyRequest struct{}
