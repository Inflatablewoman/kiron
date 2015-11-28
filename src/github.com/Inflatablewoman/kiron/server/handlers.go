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

	// Get users
	mux.Handle("GET", "/api/v1/users", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getUsers)), AuthContext{}))

	// Create user
	mux.Handle("POST", "/api/v1/users", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(createUser)), AuthContext{}))

	// Get single user
	mux.Handle("GET", "/api/v1/users/{userID}", tigertonic.WithContext(tigertonic.If(getContext, tigertonic.Marshaled(getUser)), AuthContext{}))

}

// createUser will create a user
func createUser(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *User, error) {
	var err error
	defer CatchPanic(&err, "createUser")

	log.Println("createUser Started")

	user := User{EmailAddress: "example@brainloop.com"}

	// All good!
	return http.StatusOK, nil, &user, nil
}

// getUser will get a user
func getUser(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *User, error) {
	var err error
	defer CatchPanic(&err, "getUser")

	log.Println("getUser Started")

	user := User{EmailAddress: "example@brainloop.com"}

	// All good!
	return http.StatusOK, nil, &user, nil
}

// getUsers will get a user
func getUsers(u *url.URL, h http.Header, _ interface{}) (int, http.Header, []*User, error) {
	var err error
	defer CatchPanic(&err, "getUsers")

	log.Println("getUsers Started")

	user := User{EmailAddress: "example@brainloop.com"}
	users := []*User{&user}

	// All good!
	return http.StatusOK, nil, users, nil
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
