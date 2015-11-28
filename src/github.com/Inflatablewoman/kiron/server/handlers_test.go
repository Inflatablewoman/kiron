package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {

	host := os.Getenv("KIRON_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("KIRON_PORT")
	if port == "" {
		port = "80"
	}

	curURL := fmt.Sprintf("http://%s:%s/api/v1/users", host, port)

	emailAddress := "bob@bob.com"
	firstName := "bob"
	lastName := "bobo"
	password := "bobtown"

	// Test create user
	cur := createUserRequest{EmailAddress: emailAddress, Name: firstName, LastName: lastName, Password: password}

	t.Logf("Adding user: %v", cur)

	requestBytes, err := json.Marshal(cur)
	require.NoError(t, err)

	request, err := http.NewRequest("POST", curURL, bytes.NewBuffer(requestBytes))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	response, err := client.Do(request)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	var repoUser User
	err = json.Unmarshal(body, &repoUser)
	require.NoError(t, err)

	require.Equal(t, emailAddress, repoUser.EmailAddress)
	require.Equal(t, firstName, repoUser.FirstName)
	require.Equal(t, lastName, repoUser.LastName)
	require.Equal(t, RoleAdmin, repoUser.Role)
	require.True(t, repoUser.ID > 0)
}
