package server

import (
	"fmt"
	"os"
	"testing"
	"time"

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
	cur := createUserRequest{EmailAddress: emailAddress, FirstName: firstName, LastName: lastName, Password: password}

	t.Logf("Adding user: %v", user)

	requestBytes, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", curURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var repoUser User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	require.NoError(t, err)
	require.Equal(t, emailAddress, repoUser.EmailAddress)
	require.Equal(t, firstName, repoUser.FirstName)
	require.Equal(t, lastName, repoUser.LastName)
	require.Equal(t, bcryptPassword, repoUser.Password)
	require.WithinDuration(t, created, repoUser.Created, time.Duration(5*time.Second))
	require.Equal(t, RoleAdmin, repoUser.Role)
	require.True(t, repoUser.ID > 0)
}
