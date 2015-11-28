package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {

	// Connect
	repo, err := getPostgresDB()
	require.NoError(t, err)

	emailAddress := "neil@thepetshop.boys"
	firstName := "neil"

	// Test create user
	user := User{EmailAddress: emailAddress, FirstName: firstName}

	t.Logf("Adding user: %v", user)

	err = repo.SetUser(&user)
	require.NoError(t, err)

	t.Log("User Set")

	repoUser, err := repo.GetUserByEmail(emailAddress)
	require.NoError(t, err)
	require.Equal(t, emailAddress, repoUser.EmailAddress)
	require.Equal(t, firstName, repoUser.FirstName)
	require.True(t, repoUser.ID > 0)

	t.Logf("Adding user: %v", repoUser)

	err = repo.DeleteUser(repoUser.ID)
	require.NoError(t, err)

	t.Log("User deleted")
}
