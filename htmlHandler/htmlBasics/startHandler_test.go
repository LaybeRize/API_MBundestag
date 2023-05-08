package htmlBasics

import (
	"API_MBundestag/database"
	gen "API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"regexp"
	"testing"
)

func TestStartPage(t *testing.T) {
	Setup()
	database.TestSetup()

	t.Run("testStartPageSetup", testStartPageSetup)
	t.Run("testGetStartPage", testGetStartPage)
	t.Run("testPostStartPage", testPostStartPage)
	t.Run("testLogoutPage", testLogoutPage)
}

func testLogoutPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "type=logout")

	err = acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	assert.True(t, acc.RefToken.Valid)
	assert.True(t, acc.ExpDate.Valid)
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage|testPageLogin|false|"+string(gen.SuccessfullLoggedOut)+"|true", w.Body.String())
	err = acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	assert.False(t, acc.RefToken.Valid)
	assert.False(t, acc.ExpDate.Valid)
}

func testPostStartPage(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage||false|"+string(gen.PasswordOrUsernameNotTypedIn)+"|false", w.Body.String())

	//wrong username
	w, ctx = hHa.GetContextWithForm(t, map[string]string{"username": "bazingaalsdnj", "password": "test"})
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage||false|"+string(gen.PasswordOrUsernameWrong)+"|false", w.Body.String())
	//wrong password
	w, ctx = hHa.GetContextWithForm(t, map[string]string{"username": "testPageLogin", "password": "test1"})
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage||false|"+string(gen.PasswordOrUsernameWrong)+"|false", w.Body.String())

	//test LoginTries
	acc := database.Account{}
	err := acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	acc.LoginTries = 10
	err = acc.SaveChanges()
	assert.Nil(t, err)

	w, ctx = hHa.GetContextWithForm(t, map[string]string{"username": "testPageLogin", "password": "test1"})
	PostStartPage(ctx)
	err = acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	assert.Equal(t, "TestLoginPage||false|"+acc.NextLoginTime.Time.Format(gen.FormatStringForLoginTimeout)+"|false", w.Body.String())

	w, ctx = hHa.GetContextWithForm(t, map[string]string{"username": "testPageLogin", "password": "test"})
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage||false|"+acc.NextLoginTime.Time.Format(gen.FormatStringForLoginTimeout)+"|false", w.Body.String())

	//test suspended
	err = acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	acc.LoginTries = 0
	acc.NextLoginTime.Valid = false
	acc.Suspended = true
	err = acc.SaveChanges()
	assert.Nil(t, err)

	w, ctx = hHa.GetContextWithForm(t, map[string]string{"username": "testPageLogin", "password": "test"})
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage||false|"+string(gen.AccountIsSuspended)+"|false", w.Body.String())

	//success
	err = acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	acc.Suspended = false
	err = acc.SaveChanges()
	assert.Nil(t, err)

	w, ctx = hHa.GetContextWithForm(t, map[string]string{"username": "testPageLogin", "password": "test"})
	PostStartPage(ctx)
	assert.Equal(t, "TestLoginPage|testPageLogin|true|"+string(gen.SuccessfulLoggedIn)+"|true", w.Body.String())
	str := w.Header().Values("Set-Cookie")[0]
	r := regexp.MustCompile(`(?m)token=([^;]*);`)
	str = r.FindStringSubmatch(str)[1]

	err = acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)
	assert.Equal(t, str, acc.RefToken.String)
	assert.Equal(t, true, acc.RefToken.Valid)
}

func testGetStartPage(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	GetStartPage(ctx)
	assert.Equal(t, "TestLoginPage||false||false", w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("testPageLogin")
	assert.Nil(t, err)

	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetStartPage(ctx)
	assert.Equal(t, "TestLoginPage|testPageLogin|true||false", w.Body.String())
}

func testStartPageSetup(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "testPageLogin",
		Username:    "testPageLogin",
		Password:    string(hash),
		Role:        database.User,
		Linked:      sql.NullInt64{},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	var temp *template.Template
	temp, err = template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestLoginPage|{{.Page.Account.DisplayName}}|{{.Page.LoggedIn}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "start")
}
