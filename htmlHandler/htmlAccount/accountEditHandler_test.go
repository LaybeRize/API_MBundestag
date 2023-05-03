package htmlAccount

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	wr "API_MBundestag/htmlWrapper"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html"
	"html/template"
	"testing"
)

func TestAccountEditHandler(t *testing.T) {
	database.TestSetup()

	t.Run("setupSitesEditAccount", setupSitesEditAccount)
	t.Run("testGetEditAccount", testGetEditAccount)
	t.Run("testPostQueryEditAccount", testPostQueryEditAccount)
	t.Run("testChangeEditAccount", testChangeEditAccount)

}

func testChangeEditAccount(t *testing.T) {
	accChange := database.Account{}
	err := accChange.GetByUserName("test_admin_AccEditHandler")
	assert.Nil(t, err)
	//test account does not exists message
	acc := database.Account{}
	err = acc.GetByUserName("head_admin_AccEditHandler")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "change=true")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataLogic.Account{})+"|"+string(generics.CanNotChangeNoExistentAccount)+
		"\n|false", html.UnescapeString(w.Body.String()))
	//test change only flair
	dataAcc := dataLogic.Account{}
	var msg help.Message
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"username": "head_admin_AccEditHandler", "changeFlair": "true", "flair": "test"}, "change=true")
	PostEditUserPage(ctx)
	dataAcc.GetUser("head_admin_AccEditHandler", "", &msg, &dataAcc.ChangeFlair)
	assert.Equal(t, "head_admin_AccEditHandler", dataAcc.Username)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(dataLogic.CouldChangeAccount)+
		"\n|true", html.UnescapeString(w.Body.String()))
	err = acc.GetByUserName("head_admin_AccEditHandler")
	assert.Nil(t, err)
	assert.Equal(t, "test", acc.Flair)
	//test error for root admin account
	dataAcc.GetUser("head_admin", "", &msg, &dataAcc.ChangeFlair)
	assert.Equal(t, "head_admin", dataAcc.Username)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"username": "head_admin"}, "change=true")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(generics.CanNotChangeRootAccount)+
		"\n|false", html.UnescapeString(w.Body.String()))
	//test error for other head admin account
	dataAcc.GetUser("failChangeToMe", "", &msg, &dataAcc.ChangeFlair)
	assert.Equal(t, "failChangeToMe", dataAcc.Username)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"username": "failChangeToMe"}, "change=true")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(generics.DisallowedChangeToHeadAdmin)+
		"\n|false", html.UnescapeString(w.Body.String()))
	//test error for linked value
	dataAcc.GetUser("test_admin_AccEditHandler", "", &msg, &dataAcc.ChangeFlair)
	assert.Equal(t, "test_admin_AccEditHandler", dataAcc.Username)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"username": "test_admin_AccEditHandler"}, "change=true")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(generics.LinkedValueNotANumberError)+
		"\n|false", html.UnescapeString(w.Body.String()))
	//test error for role value
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"username": "test_admin_AccEditHandler", "linked": "0"}, "change=true")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(generics.RoleCanNotBeSelectedError)+
		"\n|false", html.UnescapeString(w.Body.String()))
	//test error for linked value
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"username": "test_admin_AccEditHandler",
		"flair": "testChange", "changeFlair": "false", "suspended": "true", "removeTitles": "true", "linked": "12",
		"role": string(database.User)}, "change=true")
	PostEditUserPage(ctx)
	dataAcc.GetUser("test_admin_AccEditHandler", "", &msg, &dataAcc.ChangeFlair)
	assert.Equal(t, "test_admin_AccEditHandler", dataAcc.Username)
	assert.Equal(t, int64(0), dataAcc.Linked)
	assert.Equal(t, "", dataAcc.Flair)
	assert.Equal(t, true, dataAcc.Suspended)
	assert.Equal(t, database.User, dataAcc.Role)
	dataAcc.Linked = 12
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(dataLogic.CouldChangeAccount)+
		"\n|true", html.UnescapeString(w.Body.String()))
}

func testPostQueryEditAccount(t *testing.T) {
	accChange := database.Account{}
	err := accChange.GetByUserName("test_admin_AccEditHandler")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, accChange)
	PostEditUserPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err = acc.GetByUserName("head_admin_AccEditHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{})
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataLogic.Account{})+"|"+string(generics.InvalidType)+
		"\n|false", html.UnescapeString(w.Body.String()))

	dataAcc := dataLogic.Account{
		ID:          accChange.ID,
		DisplayName: "test_admin_AccEditHandler",
		Username:    "test_admin_AccEditHandler",
		Role:        database.Admin,
	}
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "test_admin_AccEditHandler"}, "type=user")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(dataLogic.CouldFindAccount)+
		"\n|true", html.UnescapeString(w.Body.String()))

	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "test_admin_AccEditHandler"}, "type=display")
	PostEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataAcc)+"|"+string(dataLogic.CouldFindAccount)+
		"\n|true", html.UnescapeString(w.Body.String()))
}

func testGetEditAccount(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_admin_AccEditHandler")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetEditUserPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("head_admin_AccEditHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetEditUserPage(ctx)
	assert.Equal(t, "TestEditUser|"+fmt.Sprint(dataLogic.Account{})+"||false", html.UnescapeString(w.Body.String()))
}

func setupSitesEditAccount(t *testing.T) {
	acc := database.Account{
		DisplayName: "test_admin_AccEditHandler",
		Username:    "test_admin_AccEditHandler",
		Password:    "test_admin",
		Role:        database.Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "head_admin_AccEditHandler",
		Username:    "head_admin_AccEditHandler",
		Password:    "test_admin",
		Role:        database.HeadAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "failChangeToMe",
		Username:    "failChangeToMe",
		Password:    "test_admin",
		Role:        database.HeadAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestEditUser|{{.Page.Account}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "editUser")
}
