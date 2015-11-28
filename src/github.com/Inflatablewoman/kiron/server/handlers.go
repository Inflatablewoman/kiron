package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/rcrowley/go-tigertonic"
)

// AuthContext information for Marshaled calls
type AuthContext struct {
	UserAgent  string
	RemoteAddr string
}

// getContext is used to return UserAgent and Request info from the request
func getContext(r *http.Request) (http.Header, error) {

	tigertonic.Context(r).(*AuthContext).UserAgent = r.UserAgent()
	tigertonic.Context(r).(*AuthContext).RemoteAddr = RequestAddr(r)

	// TODO :  Do some basic auth

	return nil, nil
}

// RegisterHTTPHandlers registers the http handlers
func RegisterHTTPHandlers(mux *tigertonic.TrieServeMux) {

	// Create user
	mux.Handle("POST", "/api/v1/login", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(login)), AuthContext{}))

	// Create user
	mux.Handle("POST", "/api/v1/users", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(createUser)), AuthContext{}))

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
	mux.Handle("GET", "/api/v1/users/{userID}/application/{applicationID}/documents", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getDocuments)), AuthContext{}))

	// Create documents
	mux.Handle("POST", "/api/v1/users/{userID}/application/{applicationID}/documents", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(createDocuments)), AuthContext{}))

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
func login(u *url.URL, h http.Header, request *loginRequest, context *AuthContext) (int, http.Header, *LoginResponse, error) {
	var err error
	defer CatchPanic(&err, "login")

	log.Println("login called by: %s %s", context.RemoteAddr, context.UserAgent)

	// All good!
	return http.StatusOK, nil, nil, nil
}

type createUserRequest struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	LastName     string `json:"lastname"`
}

// createUser will create a user
func createUser(u *url.URL, h http.Header, request *createUserRequest, context *AuthContext) (int, http.Header, *User, error) {
	var err error
	defer CatchPanic(&err, "createUser")

	log.Println("createUser called by: %s %s", context.RemoteAddr, context.UserAgent)

	hashedPassword, _ := createHashedPassword(request.Password)
	user := User{EmailAddress: request.EmailAddress, Password: hashedPassword, FirstName: request.Name, LastName: request.LastName, Role: RoleAdmin}

	err = repository.SetUser(&user)

	if err != nil {
		log.Printf("Unable to set user: %v", err)
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, &user, nil
}

// getUser will get a user
func getUser(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *User, error) {
	var err error
	defer CatchPanic(&err, "getUser")

	log.Println("getUser Started")

	userID := u.Query().Get("userID")

	user, err := repository.GetUser(userID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, user, nil
}

// getUsers will get a user
func getUsers(u *url.URL, h http.Header, _ interface{}) (int, http.Header, []User, error) {
	var err error
	defer CatchPanic(&err, "getUsers")

	log.Println("getUsers Started")

	users, err := repository.GetUsers()

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, users, nil
}

// getApplications will get a list of applications
func getApplications(u *url.URL, h http.Header, _ interface{}) (int, http.Header, []Application, error) {
	var err error
	defer CatchPanic(&err, "getApplications")

	log.Println("getApplications Started")

	applications, err := repository.GetApplications()

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, applications, nil
}

// getApplication will get a list of applications
func getApplication(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *Application, error) {
	var err error
	defer CatchPanic(&err, "getApplication")

	log.Println("getApplication Started")

	userID, err := strconv.Atoi(u.Query().Get("userID"))

	application, err := repository.GetApplicationOf(userID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, application, nil
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
	StudyProgram          string    `json:"study_program"`
	EducationLevel        int       `json:"education_level_id"`
}

func createApplication(u *url.URL, h http.Header, request createApplicationRequest) (int, http.Header, *Application, error) {
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
	return http.StatusCreated, nil, &application, nil
}

func getDocuments(u *url.URL, h http.Header, _ interface{}) (int, http.Header, [][]byte, error) {
	var err error
	defer CatchPanic(&err, "getDocuments")

	log.Println("getDocuments Started")

	userID := u.Query().Get("userID")
	applicationID := u.Query().Get("applicationID")

	documents, err := repository.GetDocuments(userID, applicationID)

	if err != nil {
		return http.StatusInternalServerError, nil, nil, nil
	}

	// All good!
	return http.StatusOK, nil, documents, nil
}

func createDocuments(u *url.URL, h http.Header, _ interface{}) (int, http.Header, []*Document, error) {
	var err error
	defer CatchPanic(&err, "createDocuments")

	log.Println("createDocuments Started")

	// TODO Implement functionality
	//documents := nil

	// All good!
	return http.StatusOK, nil, nil, nil
}

func getComments(u *url.URL, h http.Header, _ interface{}) (int, http.Header, []*Comment, error) {
	var err error
	defer CatchPanic(&err, "getComments")

	log.Println("getComments Started")

	// TODO Implement functionality
	//comments := nil

	// All good!
	return http.StatusOK, nil, nil, nil
}

func createComment(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *Comment, error) {
	var err error
	defer CatchPanic(&err, "createComment")

	log.Println("createComment Started")

	// TODO Implement functionality
	//comment := nil

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

	//userID := r.URL.Query().Get("userID")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleErrorWithResponse(w, err)
		return
	}

	log.Printf("Save body %v", body)

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

	userID := r.URL.Query().Get("userID")
	documentID := r.URL.Query().Get("documentID")

	log.Printf("Download User [%v] Document [%v]", userID, documentID)

	header := w.Header()
	header["Content-Type"] = []string{"application/octet-stream"}

	// TODO:  Add data to stream
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
