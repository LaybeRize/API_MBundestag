package dataLogic

import (
	"API_MBundestag/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createStandardContext(method string, URL string, w http.ResponseWriter) (*gin.Context, *gin.Engine) {
	c, r := gin.CreateTestContext(w)
	request, _ := http.NewRequest(method, URL, nil)
	c.Request = request
	return c, r
}

func TestSessionValidation(t *testing.T) {
	database.TestSetup()

	t.Run("testValidToken", testValidToken)
	t.Run("testValidTokenInvalidRole", testValidTokenInvalidRole)
	t.Run("testRandomConnection", testRandomConnection)
	t.Run("testInvalidValue", testInvalidValue)
	t.Run("testExperationDate", testExperationDate)
	t.Run("testSuspended", testSuspended)
	t.Run("testSetCookieForAccount", testSetCookieForAccount)
	t.Run("testClearCookie", testClearCookie)
}

func testClearCookie(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("Richard Falkendorf")
	assert.Nil(t, err)
	assert.Equal(t, true, acc.RefToken.Valid)
	assert.Equal(t, true, acc.ExpDate.Valid)
	err = ClearCookieFromUser(acc)
	assert.Nil(t, err)
	err = acc.GetByDisplayName("Richard Falkendorf")
	assert.Nil(t, err)
	assert.Equal(t, false, acc.RefToken.Valid)
	assert.Equal(t, false, acc.ExpDate.Valid)
}

func testSetCookieForAccount(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("Richard Falkendorf")
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	c, _ := createStandardContext("", "", w)
	err = RequestToSetCookieForAccount(acc, c)
	assert.Nil(t, err)
	val := c.Writer.Header().Get("Set-Cookie")
	assert.Nil(t, err)
	err = acc.GetByUserName(acc.Username)
	assert.Nil(t, err)
	assert.Equal(t, "token="+acc.RefToken.String+"; Path=/; Max-Age=604800", val)
}

func testSuspended(t *testing.T) {
	c, acc := testHelper(t)
	acc.Suspended = true
	err := acc.SaveChanges()
	assert.Nil(t, err)

	acc, correct := CheckUserPrivileged(c, database.User, database.MediaAdmin)
	assert.Equal(t, false, correct)
	assert.Equal(t, database.NotLoggedIn, acc.Role)

}

func testExperationDate(t *testing.T) {
	c, acc := testHelper(t)
	acc.ExpDate.Time = time.Now().Add(-1 * time.Hour)
	assert.True(t, acc.ExpDate.Time.Before(time.Now()))
	err := acc.SaveChanges()
	assert.Nil(t, err)

	var correct bool
	acc, correct = CheckUserPrivileged(c, database.User, database.PressAccount)
	assert.Equal(t, false, correct)
	assert.Equal(t, database.NotLoggedIn, acc.Role)

}

func testInvalidValue(t *testing.T) {
	c, acc := testHelper(t)
	acc.ExpDate.Valid = false
	err := acc.SaveChanges()
	assert.Nil(t, err)

	acc, correct := CheckUserPrivileged(c, database.User)
	assert.Equal(t, false, correct)
	assert.Equal(t, database.NotLoggedIn, acc.Role)
}

func testRandomConnection(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := createStandardContext("", "", w)

	acc, correct := CheckUserPrivileged(c, database.User)
	assert.Equal(t, false, correct)
	assert.Equal(t, database.NotLoggedIn, acc.Role)
}

func testValidTokenInvalidRole(t *testing.T) {
	c, acc := testHelper(t)

	acc, correct := CheckUserPrivileged(c, database.User, database.Admin)
	assert.Equal(t, false, correct)
	assert.Equal(t, database.HeadAdmin, acc.Role)
}

func testHelper(t *testing.T) (*gin.Context, database.Account) {
	account := database.Account{}
	err := account.GetByDisplayName("Richard Falkendorf")
	assert.Nil(t, err)
	account.ExpDate.Time = time.Now().UTC().Add(time.Second * time.Duration(timeUntilTokenRunsOut))
	account.ExpDate.Valid = true
	account.RefToken.String = uuid.New().String()
	account.RefToken.Valid = true
	err = account.SaveChanges()
	assert.Nil(t, err)
	err = account.GetByUserName("LaybeR")
	assert.Nil(t, err)
	cookie := http.Cookie{Name: "token", Value: account.RefToken.String, Path: "/", Domain: "localhost", MaxAge: timeUntilTokenRunsOut}
	w := httptest.NewRecorder()
	c, _ := createStandardContext("", "", w)
	c.Request.AddCookie(&cookie)
	return c, account
}

func testValidToken(t *testing.T) {
	account := database.Account{
		DisplayName: "Richard Falkendorf",
		Flair:       "",
		Username:    "LaybeR",
		Password:    "test",
		Role:        database.HeadAdmin,
	}
	err := account.CreateMe()
	assert.Nil(t, err)

	c, account := testHelper(t)

	acc, correct := CheckUserPrivileged(c, database.User, database.HeadAdmin)
	assert.Equal(t, true, correct)
	assert.NotEqual(t, account, acc)
	assert.Equal(t, account.DisplayName, acc.DisplayName)
}
