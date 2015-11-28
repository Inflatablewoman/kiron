package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {

	// Connect
	repo, err := getPostgresDB()
	require.NoError(t, err)

	emailAddress := "neil@thepetshop.boys"
	firstName := "neil"
	lastName := "waterman"

	password := "westEndGirls"
	// In test have very low complex
	bcryptPassword, err := createHashedPassword(password)
	require.NoError(t, err)

	// Test create user
	user := User{EmailAddress: emailAddress, FirstName: firstName, LastName: lastName, Password: bcryptPassword, Created: time.Now().UTC(), Role: RoleAdmin}

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
