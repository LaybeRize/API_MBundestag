package htmlBasics

import (
	"API_MBundestag/database"
	gen "API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"testing"
	"time"
)

func TestStartPage(t *testing.T) {
	Setup()

	database.TestSetup()
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "test",
		Username:    "test",
		Password:    string(hash),
		Role:        database.User,
		Linked:      sql.NullInt64{},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	gin.SetMode(gin.TestMode)

	t.Run("setupStartPage", setupStartPage)
	t.Run("testAccountLoginTries", testAccountLoginTries)
	t.Run("testGetLoggedOutStruct", testGetLoggedOutStruct)
	t.Run("testValidateUserLogin", testValidateUserLogin)
	t.Run("testGetStartPage", testGetStartPage)
	t.Run("testPostStartPage", testPostStartPage)
	t.Run("testLogoutPage", testLogoutPage)
}

func testLogoutPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("test")
	assert.Nil(t, err)
	assert.Equal(t, true, acc.RefToken.Valid)
	assert.Equal(t, true, acc.ExpDate.Time.After(time.Now()))
	w, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "type=logout"
	ctx.Request.AddCookie(&http.Cookie{
		Name:    "token",
		Value:   acc.RefToken.String,
		Expires: time.Now().Add(time.Hour),
	})
	PostStartPage(ctx)
	assert.Contains(t, "start "+gen.SuccessfullLoggedOut, w.Body.String())
	err = acc.GetByUserName("test")
	assert.Nil(t, err)
	assert.Equal(t, false, acc.RefToken.Valid)
	assert.Equal(t, false, acc.ExpDate.Valid)
}

func testPostStartPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc.LoginTries = 0
	acc.NextLoginTime.Valid = false
	err = acc.SaveChanges()
	assert.Nil(t, err)

	w, ctx := htmlHandler.GetEmptyContext(t)
	PostStartPage(ctx)
	assert.Equal(t, "start "+gen.PasswordOrUsernameNotTypedIn, w.Body.String())

	w, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"username": "test", "password": "test"})
	PostStartPage(ctx)
	assert.Equal(t, "start test", w.Body.String())
}

func testGetStartPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.AddCookie(&http.Cookie{
		Name:    "token",
		Value:   "adfasd",
		Expires: time.Now().Add(time.Hour),
	})
	GetStartPage(ctx)
	assert.Equal(t, "start ", w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("test")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)

	GetStartPage(ctx)
	assert.Contains(t, "start test", w.Body.String())

}

func testValidateUserLogin(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc.LoginTries = 0
	acc.NextLoginTime.Valid = false
	acc.NextLoginTime.Time = time.Time{}
	err = acc.SaveChanges()
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetEmptyContext(t)
	get := validateUserLogin(ctx)
	assert.Equal(t, gen.PasswordOrUsernameNotTypedIn, get.Info)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "bruh", "username": "asdf"})
	get = validateUserLogin(ctx)
	assert.Equal(t, gen.PasswordOrUsernameWrong, get.Info)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "test", "username": "test"})
	get = validateUserLogin(ctx)
	assert.Equal(t, "", get.Info)
	assert.Equal(t, acc, get.Account)

	acc.Suspended = true
	err = acc.SaveChanges()
	assert.Nil(t, err)
	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "test", "username": "test"})
	get = validateUserLogin(ctx)
	assert.Equal(t, gen.AccountIsSuspended, get.Info)
	acc.Suspended = false
	err = acc.SaveChanges()
	assert.Nil(t, err)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "testa", "username": "test"})
	get = validateUserLogin(ctx)
	assert.Equal(t, gen.PasswordOrUsernameWrong, get.Info)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "testa", "username": "test"})
	get = validateUserLogin(ctx)
	assert.Equal(t, gen.PasswordOrUsernameWrong, get.Info)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "testa", "username": "test"})
	get = validateUserLogin(ctx)
	assert.Equal(t, gen.PasswordOrUsernameWrong, get.Info)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "testa", "username": "test"})
	get = validateUserLogin(ctx)
	err = acc.GetByUserName("test")
	assert.Nil(t, err)
	assert.Equal(t, acc.NextLoginTime.Time.Format(gen.FormatStringForLoginTimeout), get.Info)

	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"password": "testa", "username": "test"})
	get = validateUserLogin(ctx)
	assert.Equal(t, acc.NextLoginTime.Time.Format(gen.FormatStringForLoginTimeout), get.Info)
}

func testGetLoggedOutStruct(t *testing.T) {
	get := getLoggedOutStartPageStruct("test")
	assert.Equal(t, StartPageStruct{
		Account:  database.Account{Role: database.NotLoggedIn},
		LoggedIn: false,
		Info:     "test",
	}, *get)
}

func testAccountLoginTries(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("test")
	assert.Nil(t, err)
	acc.LoginTries = 2
	err = updateLoginTries(&acc)

	assert.Nil(t, err)
	assert.Equal(t, 3, acc.LoginTries)
	assert.Equal(t, false, acc.NextLoginTime.Valid)

	acc.LoginTries = 3
	err = updateLoginTries(&acc)
	assert.Equal(t, AccountCanNotBeLoggindBecauseOfTimeout, err)

	assert.Equal(t, 4, acc.LoginTries)
	assert.Equal(t, true, acc.NextLoginTime.Valid)
}

func setupStartPage(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{if .Page.LoggedIn}}{{.Page.Account.DisplayName}}{{end}}{{.Page.Info}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"start": temp,
		},
	}
}
