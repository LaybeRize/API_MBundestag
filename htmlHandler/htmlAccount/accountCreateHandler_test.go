package htmlAccount

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"html"
	"html/template"
	"strconv"
	"testing"
)

func TestAccountCreateHandler(t *testing.T) {
	database.TestSetup()
	Setup()
	htmlBasics.Setup()

	t.Run("setupSitesCreateAccount", setupSitesCreateAccount)
	t.Run("testGetCreateAccount", testGetCreateAccount)
	t.Run("testPostCreateAccount", testPostCreateAccount)
}

func testPostCreateAccount(t *testing.T) {
	//error Test
	acc := database.Account{}
	err := acc.GetByUserName("test_admin_AccCreateHandler")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostCreateUserPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	//test fail linked value
	err = acc.GetByUserName("head_admin_AccCreateHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	PostCreateUserPage(ctx)
	assert.Equal(t, "TestCreateUser|"+fmt.Sprint(database.Account{Role: database.User})+"|"+
		string(generics.LinkedValueNotANumberError)+"|false", html.UnescapeString(w.Body.String()))
	//test fail role value
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"linked": "12"})
	PostCreateUserPage(ctx)
	assert.Equal(t, "TestCreateUser|"+fmt.Sprint(database.Account{Role: database.User, Linked: sql.NullInt64{Int64: 12}})+"|"+
		string(generics.RoleCanNotBeSelectedError)+"|false", html.UnescapeString(w.Body.String()))
	//test password value
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"linked": "12", "role": string(database.PressAccount)})
	PostCreateUserPage(ctx)
	assert.Equal(t, "TestCreateUser|"+fmt.Sprint(database.Account{Role: database.PressAccount,
		Linked: sql.NullInt64{Valid: true, Int64: 12}})+"|"+
		string(generics.NamesOrPasswordIsEmptyError)+"|false", html.UnescapeString(w.Body.String()))

	//success create User
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"linked": "0", "role": string(database.User),
		"displayname": "test_CreationHTML", "username": "test_CreationHTML", "password": "test"})
	PostCreateUserPage(ctx)
	newAcc := database.Account{}
	err = newAcc.GetByDisplayName("test_CreationHTML")
	assert.Nil(t, err)
	assert.Equal(t, "TestCreateUser|"+fmt.Sprint(database.Account{Role: database.User,
		Username: "test_CreationHTML", DisplayName: "test_CreationHTML", ID: newAcc.ID})+"|"+
		string(generics.SuccesFullCreatedAccount)+"|true", html.UnescapeString(w.Body.String()))
	assert.Equal(t, "test_CreationHTML", newAcc.DisplayName)
	assert.Equal(t, "test_CreationHTML", newAcc.Username)
	assert.Equal(t, database.User, newAcc.Role)
	err = bcrypt.CompareHashAndPassword([]byte(newAcc.Password), []byte("test"))
	assert.Nil(t, err)

	//success create PressAccount
	accID := newAcc.ID
	accString := strconv.FormatInt(newAcc.ID, 10)
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"linked": accString, "role": string(database.PressAccount),
		"displayname": "test_CreationPressAccount"})
	PostCreateUserPage(ctx)
	err = newAcc.GetByDisplayName("test_CreationPressAccount")
	assert.Nil(t, err)
	assert.Equal(t, "TestCreateUser|"+fmt.Sprint(database.Account{Role: database.PressAccount,
		Username: "test_CreationPressAccount", DisplayName: "test_CreationPressAccount", ID: newAcc.ID,
		Linked: sql.NullInt64{Valid: true, Int64: accID}})+"|"+
		string(generics.SuccesFullCreatedAccount)+"|true", html.UnescapeString(w.Body.String()))
	assert.Equal(t, "test_CreationPressAccount", newAcc.DisplayName)
	assert.Equal(t, "test_CreationPressAccount", newAcc.Username)
	assert.Equal(t, database.PressAccount, newAcc.Role)
	assert.Equal(t, "", newAcc.Password)
	assert.Equal(t, sql.NullInt64{Valid: true, Int64: accID}, newAcc.Linked)
}

func testGetCreateAccount(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_admin_AccCreateHandler")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetCreateUserPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("head_admin_AccCreateHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetCreateUserPage(ctx)
	assert.Equal(t, "TestCreateUser|"+fmt.Sprint(database.Account{Role: database.User})+"||false", html.UnescapeString(w.Body.String()))
}

func setupSitesCreateAccount(t *testing.T) {
	acc := database.Account{
		DisplayName: "test_admin_AccCreateHandler",
		Username:    "test_admin_AccCreateHandler",
		Password:    "test_admin",
		Role:        database.Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "head_admin_AccCreateHandler",
		Username:    "head_admin_AccCreateHandler",
		Password:    "test_admin",
		Role:        database.HeadAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestCreateUser|{{.Page.Account}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "createUser")
}
