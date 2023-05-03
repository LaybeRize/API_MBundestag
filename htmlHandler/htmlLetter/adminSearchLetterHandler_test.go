package htmlLetter

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"testing"
	"time"
)

func TestAdminSearchLetterHandler(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestLettersDB()

	t.Run("setupAccountsAndLettersForAdminSearch", setupAccountsAndLettersForAdminSearch)
	t.Run("setupTestAdminSearchLetterPage", setupTestAdminSearchLetterPage)
	t.Run("testGetAdminSearchLetterPage", testGetAdminSearchLetterPage)
	t.Run("testPostAdminSearchLetterPage", testPostAdminSearchLetterPage)
}

func testPostAdminSearchLetterPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostAdminLetterViewPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("testAdmin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"uuid": "lol"})
	PostAdminLetterViewPage(ctx)
	assert.Equal(t, "adminViewLetter lol "+generics.ErrorUUIDDoesNotExist, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"uuid": "testUUID"})
	PostAdminLetterViewPage(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/admin-letter-view?uuid=testUUID", w.Header().Get("Location"))
}

func testGetAdminSearchLetterPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetAdminLetterViewPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("testAdmin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=lol"
	GetAdminLetterViewPage(ctx)
	assert.Equal(t, "adminViewLetter lol "+generics.ErrorUUIDDoesNotExist, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=testUUID"
	GetAdminLetterViewPage(ctx)
	assert.Equal(t, "viewLetter lol  false [user1 user2] [user1] [user2]", w.Body.String())

}

func setupTestAdminSearchLetterPage(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Letter.Title}} {{.Page.Account.DisplayName}} {{.Page.Letter.Info.ModMessage}} {{.Page.Letter.Info.PeopleInvitedToSign}} {{.Page.Letter.Info.PeopleNotYetSigned}} {{.Page.Letter.Info.Signed}}")
	assert.Nil(t, err)
	temp3, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.UUID}} {{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":           temp,
			"viewLetter":      temp2,
			"adminViewLetter": temp3,
		},
	}
}

func setupAccountsAndLettersForAdminSearch(t *testing.T) {
	htmlHandler.CreateAccountForTest(t, "testAdmin", database.Admin, 0)
	htmlHandler.CreateAccountForTest(t, "user1", database.User, 0)
	htmlHandler.CreateAccountForTest(t, "user2", database.User, 0)
	letter := database.Letter{
		UUID:        "testUUID",
		Written:     time.Time{},
		Author:      "Nutzer 1",
		Flair:       "TestFlair",
		Title:       "lol",
		Content:     "bruh",
		HTMLContent: "bazinga",
		Info: database.LetterInfo{
			ModMessage:          false,
			AllHaveToAgree:      false,
			NoSigning:           false,
			PeopleInvitedToSign: []string{"user1", "user2"},
			PeopleNotYetSigned:  []string{"user1"},
			Signed:              []string{"user2"},
			Rejected:            []string{},
		},
	}
	err := letter.CreateMe()
	assert.Nil(t, err)
}
