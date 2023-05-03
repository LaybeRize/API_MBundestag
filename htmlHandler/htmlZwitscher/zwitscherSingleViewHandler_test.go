package htmlZwitscher

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
	"time"
)

func TestZwitscherSingleViewPage(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestZwitscherDB()
	database.TestAccountDB()

	t.Run("setupTestZwitscherSingleViewPage", setupTestZwitscherSingleViewPage)
	t.Run("testNotAuthorizedZwitscherSingle", testNotAuthorizedZwitscherSingle)
	t.Run("testGetZwitscherSingleViewPage", testGetZwitscherSingleViewPage)
	t.Run("testPostZwitscherSingleViewPage", testPostZwitscherSingleViewPage)
}

func testPostZwitscherSingleViewPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "admin", "content": "Test"})
	ctx.Request.URL.RawQuery = "uuid=test2"
	PostZwitscherLatestViewPage(ctx)
	assert.Equal(t, "viewZwitscher Test Author TestCotent Test Author2 TestContent2 true true admin  "+generics.ZwitscherCreationSuccessful+"\n admin Test|", w.Body.String())
}

func testGetZwitscherSingleViewPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test2"
	GetZwitscherLatestViewPage(ctx)
	assert.Equal(t, "viewZwitscher Test Author "+generics.ZwitscherBlockText+" Test Author2 TestContent2 false false    ", w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test2"
	GetZwitscherLatestViewPage(ctx)
	assert.Equal(t, "viewZwitscher Test Author TestCotent Test Author2 TestContent2 true true    ", w.Body.String())
}

func testNotAuthorizedZwitscherSingle(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test2"
	PostZwitscherLatestViewPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
}

func setupTestZwitscherSingleViewPage(t *testing.T) {
	exmp := database.Zwitscher{
		UUID:        "test",
		Author:      "Test Author",
		HTMLContent: "TestCotent",
		Blocked:     true,
	}
	err := exmp.CreateMe()
	assert.Nil(t, err)
	err = exmp.SaveChanges()
	assert.Nil(t, err)
	time.Sleep(100)
	exmp.UUID, exmp.Author, exmp.HTMLContent, exmp.Blocked = "test2", "Test Author2", "TestContent2", false
	exmp.ConnectedTo = sql.NullString{Valid: true, String: "test"}
	err = exmp.CreateMe()
	assert.Nil(t, err)
	htmlHandler.CreateAccountForTest(t, "admin", database.MediaAdmin, 0)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Parent.Author}} {{.Page.Parent.HTMLContent}} {{.Page.Self.Author}} {{.Page.Self.HTMLContent}} {{.Page.CanZwitscher}} {{.Page.CanSuspendZwitscher}} {{.Page.SelectedAccount}} {{.Page.Content}} {{.Page.Message}} {{range $i, $pub := .Page.Zwitscher}}{{$pub.Author}} {{$pub.HTMLContent}}|{{end}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":         temp,
			"viewZwitscher": temp2,
		},
	}
}

func TestZwitscherSingleStructTest(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestZwitscherDB()
	database.TestAccountDB()

	htmlHandler.CreateAccountForTest(t, "user", database.User, 0)
	t.Run("testWriteAccountDoesNotExistInDatabaseZwitscherSingle", testWriteAccountDoesNotExistInDatabaseZwitscherSingle)
	t.Run("testThatAccountIsNotAllowedToWriteZwitscherSingle", testThatAccountIsNotAllowedToWriteZwitscherSingle)
	t.Run("testNotEmptyZwitscherSingle", testNotEmptyZwitscherSingle)
	t.Run("testZwitscherIsToLongZwitscherSingle", testZwitscherIsToLongZwitscherSingle)
	t.Run("testZwitscherCreationSuccessfulSingle", testZwitscherCreationSuccessfulSingle)
}

func testZwitscherCreationSuccessfulSingle(t *testing.T) {
	arra := database.ZwitscherList{}
	err := arra.GetCommentsFor("test", true)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(arra))

	res := ZwitscherSingleViewStruct{}
	acc := database.Account{}
	err = acc.GetByUserName("user")
	assert.Nil(t, err)
	res.Self.UUID = "test"
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user", "content": "Test"})
	res.validateMakeComment(ctx, &database.Account{DisplayName: "user", ID: 1})

	assert.Equal(t, generics.ZwitscherCreationSuccessful+"\n", res.Message)

	err = arra.GetCommentsFor("test", true)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(arra))
	assert.Equal(t, "Test", arra[0].HTMLContent)
}

func testZwitscherIsToLongZwitscherSingle(t *testing.T) {
	res := ZwitscherSingleViewStruct{}
	acc := database.Account{}
	err := acc.GetByUserName("user")
	assert.Nil(t, err)
	cont := htmlHandler.PadLeft("", "a", generics.CharacterLimitZwitscher+10)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user", "content": cont})
	res.validateMakeComment(ctx, &database.Account{DisplayName: "user", ID: 1})

	assert.Equal(t, fmt.Sprintf(generics.ZwitscherIsToLong, generics.CharacterLimitZwitscher)+"\n", res.Message)
}

func testNotEmptyZwitscherSingle(t *testing.T) {
	res := ZwitscherSingleViewStruct{}
	acc := database.Account{}
	err := acc.GetByUserName("user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user"})
	res.validateMakeComment(ctx, &database.Account{DisplayName: "user", ID: 1})

	assert.Equal(t, generics.ZwitscherIsNotAllowedToBeEmpty+"\n", res.Message)
}

func testThatAccountIsNotAllowedToWriteZwitscherSingle(t *testing.T) {
	res := ZwitscherSingleViewStruct{}
	acc := database.Account{}
	err := acc.GetByUserName("user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user"})
	res.validateMakeComment(ctx, &database.Account{DisplayName: "lol", ID: 12})

	assert.Equal(t, generics.AccountIsNotYours+"\n", res.Message)
}

func testWriteAccountDoesNotExistInDatabaseZwitscherSingle(t *testing.T) {
	res := ZwitscherSingleViewStruct{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	res.validateMakeComment(ctx, &database.Account{})

	assert.Equal(t, generics.AccountDoesNotExists+"\n", res.Message)
}
