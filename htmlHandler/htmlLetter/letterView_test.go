package htmlLetter

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestLetterView(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestLettersDB()
	acc := database.Account{
		DisplayName: "test",
		Username:    "test",
		Password:    "test",
		Role:        database.User,
		Linked:      sql.NullInt64{},
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.Username, acc.DisplayName = "test2", "test2"
	err = acc.CreateMe()
	assert.Nil(t, err)
	letter := database.Letter{
		UUID:        "asvasd",
		Author:      "gsdf",
		Flair:       "bcvxd",
		Title:       "sfdsdf",
		Content:     "awds",
		HTMLContent: "asd",
		Info: database.LetterInfo{
			ModMessage:          false,
			AllHaveToAgree:      true,
			NoSigning:           false,
			PeopleInvitedToSign: []string{"test"},
			PeopleNotYetSigned:  []string{"test"},
			Signed:              []string{},
			Rejected:            []string{},
		},
	}
	err = letter.CreateMe()
	assert.Nil(t, err)
	gin.SetMode(gin.TestMode)

	t.Run("setupPagesLetterView", setupPagesLetterView)
	t.Run("testFailedRequest", testFailedRequest)
	t.Run("testNotEligibleAccountRequest", testNotEligibleAccountRequest)
	t.Run("testLetterDoesNotExistRequest", testLetterDoesNotExistRequest)
	t.Run("testLetterAccountNotAllowedToViewRequest", testLetterAccountNotAllowedToViewRequest)
	t.Run("testGetLetterRequest", testGetLetterRequest)
	t.Run("testSignLetterRequest", testSignLetterRequest)
	t.Run("testRejectLetterRequest", testRejectLetterRequest)
}

func testRejectLetterRequest(t *testing.T) {
	letter := database.Letter{}
	err := letter.GetByID("asvasd")
	assert.Nil(t, err)
	letter.Info.PeopleNotYetSigned, letter.Info.Signed = []string{"test"}, []string{}
	err = letter.SaveChanges()
	assert.Nil(t, err)

	acc := database.Account{}
	err = acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=test&uuid=asvasd&type=reject"
	GetViewSingleLetter(ctx)
	assert.Equal(t, "viewLetter sfdsdf test false", w.Body.String())

	err = letter.GetByID("asvasd")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(letter.Info.PeopleNotYetSigned))
	assert.Equal(t, 1, len(letter.Info.Rejected))
	assert.Equal(t, "test", letter.Info.Rejected[0])
}

func testSignLetterRequest(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=test&uuid=asvasd&type=sign"
	GetViewSingleLetter(ctx)
	assert.Equal(t, "viewLetter sfdsdf test false", w.Body.String())

	letter := database.Letter{}
	err = letter.GetByID("asvasd")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(letter.Info.PeopleNotYetSigned))
	assert.Equal(t, 1, len(letter.Info.Signed))
	assert.Equal(t, "test", letter.Info.Signed[0])
}

func testGetLetterRequest(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=test&uuid=asvasd"
	GetViewSingleLetter(ctx)
	assert.Equal(t, "viewLetter sfdsdf test false", w.Body.String())

}

func testLetterAccountNotAllowedToViewRequest(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test2")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=test2&uuid=asvasd"
	GetViewSingleLetter(ctx)
	assert.Equal(t, "error "+generics.LetterDoesntExistOrNotAccessable, w.Body.String())
}

func testLetterDoesNotExistRequest(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	GetViewSingleLetter(ctx)
	assert.Equal(t, "error "+generics.LetterDoesntExistOrNotAccessable, w.Body.String())

}

func testNotEligibleAccountRequest(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=slkdjas"
	GetViewSingleLetter(ctx)
	assert.Equal(t, "error "+generics.AccountForLetterViewError, w.Body.String())
}

func testFailedRequest(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetViewSingleLetter(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
}

func setupPagesLetterView(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Letter.Title}} {{.Page.Account.DisplayName}} {{.Page.Letter.Info.ModMessage}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":      temp,
			"viewLetter": temp2,
		},
	}
}
