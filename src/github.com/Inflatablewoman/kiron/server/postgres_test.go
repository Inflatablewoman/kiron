package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {

	// Connect
	repo, err := getPostgresDB()
	require.NoError(t, err)

	emailAddress := fmt.Sprintf("test_%s@%s.com", GetRandomString(5, ""), GetRandomString(5, ""))
	firstName := "neil"
	lastName := "waterman"
	created := time.Now().UTC()
	password := "westEndGirls"
	// In test have very low complex
	bcryptPassword, err := createHashedPassword(password)
	require.NoError(t, err)

	// Test User functions
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

	lastName = "waterboy"
	repoUser.LastName = lastName
	err = repo.UpdateUser(repoUser)
	require.NoError(t, err)

	t.Logf("Updated User: %v", repoUser)

	repoUser, err = repo.GetUser(repoUser.ID)
	require.NoError(t, err)
	require.Equal(t, emailAddress, repoUser.EmailAddress)
	require.Equal(t, firstName, repoUser.FirstName)
	require.Equal(t, lastName, repoUser.LastName)
	require.Equal(t, bcryptPassword, repoUser.Password)
	require.WithinDuration(t, created, repoUser.Created, time.Duration(5*time.Second))
	require.Equal(t, RoleAdmin, repoUser.Role)
	require.True(t, repoUser.ID > 0)

	t.Logf("Got by id user: %v", repoUser)

	err = repo.DeleteUser(repoUser.ID)
	require.NoError(t, err)

	t.Log("User deleted")

}

func TestPostgresApplications(t *testing.T) {

	// Connect
	repo, err := getPostgresDB()
	require.NoError(t, err)

	// test Application functions
	created := time.Now().UTC()

	appl := Application{
		Birthday:              created,
		PhoneNumber:           "555",
		Nationality:           "marsian",
		Country:               "for old men",
		City:                  "atlantis",
		Zip:                   "666",
		Address:               "none",
		AddressExtra:          "of yo business",
		FirstPageOfSurveyData: "I use a GameBoy",
		Gender:                "female",
		UserID:                600,
		EducationLevel:        2,
		Status:                "rejected",
		BlockExpires:          created,
		Created:               created,
		Edited:                created}

	err = repo.SetApplication(&appl)
	require.NoError(t, err)

	repoAppl, err := repo.GetApplicationOf(600)
	require.NoError(t, err)

	t.Logf("Set application: %v", repoAppl)

	require.WithinDuration(t, created, repoAppl.Birthday, time.Duration(5*time.Second))
	require.Equal(t, "555", repoAppl.PhoneNumber)
	require.Equal(t, "for old men", repoAppl.Country)
	require.Equal(t, "marsian", repoAppl.Nationality)
	require.Equal(t, "none", repoAppl.Address)
	require.Equal(t, "of your business", repoAppl.AddressExtra)
	require.Equal(t, "atlantis", repoAppl.City)
	require.Equal(t, "female", repoAppl.Gender)
	require.Equal(t, "I use a GameBoy", repoAppl.FirstPageOfSurveyData)
	require.Equal(t, "rejected", repoAppl.Status)
	require.Equal(t, "female", repoAppl.Gender)
	require.WithinDuration(t, created, repoAppl.Created, time.Duration(5*time.Second))
	require.WithinDuration(t, created, repoAppl.Edited, time.Duration(5*time.Second))
	require.WithinDuration(t, created, repoAppl.BlockExpires, time.Duration(5*time.Second))
	require.True(t, repoAppl.ID > 0)

	t.Logf("Got Application: %v", repoAppl)

	repoAppl.Nationality = "venusian"
	err = repo.UpdateApplication(repoAppl)
	require.NoError(t, err)

	t.Logf("Updated Application: %v", repoAppl)

	err = repo.DeleteApplication(repoAppl.ID)
	require.NoError(t, err)

	t.Log("Test Application deleted")
}

func TestPostgresTokens(t *testing.T) {
	// Connect
	repo, err := getPostgresDB()
	require.NoError(t, err)

	expiry := time.Now().UTC().Add(-time.Duration(2 * time.Hour))

	token := Token{UserID: 1, Value: "Myawesometoken", Expires: expiry}

	err = repo.SetToken(&token)
	require.NoError(t, err)

	t.Logf("Set token: %v", token)

	repoToken, err := repo.GetToken(token.Value)
	require.NoError(t, err)

	require.WithinDuration(t, expiry, repoToken.Expires, time.Duration(5*time.Second))
	require.Equal(t, "Myawesometoken", repoToken.Value)
	require.Equal(t, 1, repoToken.UserID)

	t.Logf("Got Token: %v", repoToken)

	err = repo.DelToken(token.Value)
	require.NoError(t, err)

	repoToken, err = repo.GetToken(token.Value)
	// It should have been deleted
	require.Nil(t, repoToken)

	t.Log("Test Token deleted")

	// Expired already
	expiry = time.Now().UTC().Add(-time.Duration(2 * time.Hour))
	token = Token{UserID: 1, Value: "Myawesometoken", Expires: expiry}

	err = repo.SetToken(&token)
	require.NoError(t, err)

	err = repo.DelExpiredTokens()
	require.NoError(t, err)

	repoToken, err = repo.GetToken(token.Value)
	// It should have been deleted as it is an expired token
	require.Nil(t, repoToken)
}
