package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var client = http.DefaultClient

func TestCreateAndLoginUser(t *testing.T) {

	host := os.Getenv("KIRON_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("KIRON_PORT")
	if port == "" {
		port = "80"
	}

	// Create an admin in the database...
	repo, err := getPostgresDB()
	require.NoError(t, err)

	emailAddress := fmt.Sprintf("test_admin_%s@%s.com", GetRandomString(5, ""), GetRandomString(5, ""))
	firstName := "Admin"
	lastName := "MrAdmin"
	created := time.Now().UTC()
	password := "monkey"
	t.Logf("Password: '%s'", password)
	bcryptPassword, err := createHashedPassword(password)
	t.Logf("hashPassword: %s", bcryptPassword)
	require.NoError(t, err)
	user := User{EmailAddress: emailAddress, FirstName: firstName, LastName: lastName, Password: bcryptPassword, Created: created, Role: RoleAdmin}

	t.Logf("Adding user: %v", user)

	err = repo.SetUser(&user)
	require.NoError(t, err)

	// Admin user logins in... via REST
	lr := loginRequest{EmailAddress: emailAddress, Password: password}

	t.Logf("Logging in user: %v", lr)

	requestBytes, err := json.Marshal(lr)
	require.NoError(t, err)

	lurURL := fmt.Sprintf("http://%s:%s/api/v1/login", host, port)

	request, err := http.NewRequest("POST", lurURL, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	require.NoError(t, err)

	require.Equal(t, emailAddress, loginResp.Result.EmailAddress)
	require.Equal(t, firstName, loginResp.Result.FirstName)
	require.Equal(t, lastName, loginResp.Result.LastName)
	require.True(t, loginResp.Result.ID > 0)
	require.NotEmpty(t, loginResp.Token)
	require.Len(t, loginResp.Token, 16)
	require.Equal(t, loginResp.TokenExpiry, 3600)

	curURL := fmt.Sprintf("http://%s:%s/api/v1/users", host, port)

	emailAddress = fmt.Sprintf("test_user_%s@%s.com", GetRandomString(5, ""), GetRandomString(5, ""))
	firstName = "bob"
	lastName = "bobo"
	password = "bobtown"

	// Test create user
	cur := createUserRequest{EmailAddress: emailAddress, Name: firstName, LastName: lastName, Password: password}

	t.Logf("Adding user: %v", cur)

	requestBytes, err = json.Marshal(cur)
	require.NoError(t, err)

	request, err = http.NewRequest("POST", curURL, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+loginResp.Token)

	response, err = client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err = ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	var repoUser RestUser
	err = json.Unmarshal(body, &repoUser)
	require.NoError(t, err)

	require.Equal(t, emailAddress, repoUser.EmailAddress)
	require.Equal(t, firstName, repoUser.FirstName)
	require.Equal(t, lastName, repoUser.LastName)

	lr = loginRequest{EmailAddress: emailAddress, Password: password}

	t.Logf("Logging in user: %v", lr)

	requestBytes, err = json.Marshal(lr)
	require.NoError(t, err)

	lurURL = fmt.Sprintf("http://%s:%s/api/v1/login", host, port)

	request, err = http.NewRequest("POST", lurURL, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	response, err = client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err = ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &loginResp)
	require.NoError(t, err)

	require.Equal(t, emailAddress, loginResp.Result.EmailAddress)
	require.Equal(t, firstName, loginResp.Result.FirstName)
	require.Equal(t, lastName, loginResp.Result.LastName)
	require.True(t, loginResp.Result.ID > 0)
	require.NotEmpty(t, loginResp.Token)
	require.Len(t, loginResp.Token, 16)
	require.Equal(t, loginResp.TokenExpiry, 3600)

	// Create application
	// test Application functions
	created = time.Now().UTC()

	createAppReq := createApplicationRequest{
		Birthday:              created,
		PhoneNumber:           "555",
		Nationality:           "marsian",
		Country:               "for old men",
		City:                  "atlantis",
		Zip:                   "666",
		Address:               "suck it",
		AddressExtra:          "of yo business",
		FirstPageOfSurveyData: "I use a GameBoy",
		Gender:                "female",
		EducationLevel:        2,
	}

	t.Logf("Adding user: %v", createAppReq)

	requestBytes, err = json.Marshal(createAppReq)
	require.NoError(t, err)

	cAppURL := fmt.Sprintf("http://%s:%s/api/v1/users/%v/application", host, port, repoUser.ID)

	request, err = http.NewRequest("POST", cAppURL, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+loginResp.Token)

	response, err = client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err = ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	var restApplication RestApplication
	err = json.Unmarshal(body, &restApplication)
	require.NoError(t, err)

	require.NotNil(t, restApplication)
	require.Equal(t, repoUser.ID, restApplication.UserID)

	// Upload
	/*randomContent := GetRandomString(50, "test")


	// upload url
	uploadURL :=

	request, err = http.NewRequest("POST", curURL, bytes.NewBuffer([]byte(randomContent)))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+loginResp.Token)

	response, err = client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	// Download

	lurURL = fmt.Sprintf("http://%s:%s/api/v1/logout", host, port)

	er := emptyRequest{}
	requestBytes, err = json.Marshal(er)
	require.NoError(t, err)

	request, err = http.NewRequest("POST", lurURL, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+loginResp.Token)

	response, err = client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)*/

}
