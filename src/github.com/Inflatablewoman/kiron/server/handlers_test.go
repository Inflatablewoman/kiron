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

	response, err = client.Do(request)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err = ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	var repoUser User
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

}
