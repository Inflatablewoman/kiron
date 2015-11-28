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

	emailAddress := "neil2@thepetshop.boys"
	firstName := "neil"
	lastName := "waterman"
	created := time.Now().UTC()
	password := "westEndGirls"
	// In test have very low complex
	bcryptPassword, err := createHashedPassword(password)
	require.NoError(t, err)

	// Test create user
	user := User{EmailAddress: emailAddress, FirstName: firstName, LastName: lastName, Password: bcryptPassword, Created: created, Role: RoleAdmin}

	t.Logf("Adding user: %v", user)

	err = repo.SetUser(&user)
	require.NoError(t, err)

	t.Log("User Set")

	repoUser, err := repo.GetUserByEmail(emailAddress)
	require.NoError(t, err)
	require.Equal(t, emailAddress, repoUser.EmailAddress)
	require.Equal(t, firstName, repoUser.FirstName)
	require.Equal(t, lastName, repoUser.LastName)
	require.Equal(t, bcryptPassword, repoUser.Password)
	require.WithinDuration(t, created, repoUser.Created, time.Duration(5*time.Second))
	require.Equal(t, RoleAdmin, repoUser.Role)
	require.True(t, repoUser.ID > 0)

	t.Logf("Got user: %v", repoUser)

	err = repo.DeleteUser(repoUser.ID)
	require.NoError(t, err)

	t.Log("User deleted")
}
