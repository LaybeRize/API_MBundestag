package htmlZwitscher

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
	"time"
)

func TestZwitscherLeastPage(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestZwitscherDB()
	database.TestAccountDB()

	t.Run("testSetupPageZwitscherLeast", testSetupPageZwitscherLeast)
	t.Run("testNotAuthorizedZwitscherLeast", testNotAuthorizedZwitscherLeast)
	t.Run("testGetZwitscherLatestViewPage", testGetZwitscherLatestViewPage)
	t.Run("testPostZwitscherLatestViewPage", testPostZwitscherLatestViewPage)
}

func testPostZwitscherLatestViewPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "admin", "content": "Test"})
	PostZwitscherLatestViewPage(ctx)
	assert.Equal(t, "zwitscherList true admin  "+generics.ZwitscherCreationSuccessful+"\n 20 admin Test|Test Author2 TestContent2|Test Author TestCotent|", w.Body.String())
}

func testGetZwitscherLatestViewPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetZwitscherLatestViewPage(ctx)
	assert.Equal(t, "zwitscherList false    20 Test Author TestCotent|", w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetZwitscherLatestViewPage(ctx)
	assert.Equal(t, "zwitscherList true    20 Test Author2 TestContent2|Test Author TestCotent|", w.Body.String())
}

func testNotAuthorizedZwitscherLeast(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostZwitscherLatestViewPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
}

func testSetupPageZwitscherLeast(t *testing.T) {
	exmp := database.Zwitscher{
		UUID:        "test",
		Author:      "Test Author",
		HTMLContent: "TestCotent",
	}
	err := exmp.CreateMe()
	assert.Nil(t, err)
	time.Sleep(100)
	exmp.UUID, exmp.Author, exmp.HTMLContent, exmp.Blocked = "test2", "Test Author2", "TestContent2", true
	err = exmp.CreateMe()
	assert.Nil(t, err)
	err = exmp.SaveChanges()
	assert.Nil(t, err)
	htmlHandler.CreateAccountForTest(t, "admin", database.MediaAdmin, 0)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.CanZwitscher}} {{.Page.SelectedAccount}} {{.Page.Content}} {{.Page.Message}} {{.Page.Amount}} {{range $i, $pub := .Page.Zwitscher}}{{$pub.Author}} {{$pub.HTMLContent}}|{{end}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":         temp,
			"zwitscherList": temp2,
		},
	}
}

func TestZwitscherLeastStructTest(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestZwitscherDB()
	database.TestAccountDB()

	htmlHandler.CreateAccountForTest(t, "user", database.User, 0)
	t.Run("testWriteAccountDoesNotExistInDatabaseZwitscher", testWriteAccountDoesNotExistInDatabaseZwitscher)
	t.Run("testThatAccountIsNotAllowedToWriteZwitscher", testThatAccountIsNotAllowedToWriteZwitscher)
	t.Run("testNotEmptyZwitscher", testNotEmptyZwitscher)
	t.Run("testZwitscherIsToLongZwitscher", testZwitscherIsToLongZwitscher)
	t.Run("testZwitscherCreationSuccessful", testZwitscherCreationSuccessful)
}

func testZwitscherCreationSuccessful(t *testing.T) {
	arra := database.ZwitscherList{}
	err := arra.GetLatested(10, true)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(arra))

	res := ZwitscherListViewStruct{}
	acc := database.Account{}
	err = acc.GetByUserName("user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user", "content": "Test"})
	res.validateZwitscherCreate(ctx, &database.Account{DisplayName: "user", ID: 1})

	assert.Equal(t, generics.ZwitscherCreationSuccessful+"\n", res.Message)

	err = arra.GetLatested(10, true)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(arra))
	assert.Equal(t, "Test", arra[0].HTMLContent)
}

func testZwitscherIsToLongZwitscher(t *testing.T) {
	res := ZwitscherListViewStruct{}
	acc := database.Account{}
	err := acc.GetByUserName("user")
	assert.Nil(t, err)
	cont := htmlHandler.PadLeft("", "a", generics.CharacterLimitZwitscher+10)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user", "content": cont})
	res.validateZwitscherCreate(ctx, &database.Account{DisplayName: "user", ID: 1})

	assert.Equal(t, fmt.Sprintf(generics.ZwitscherIsToLong, generics.CharacterLimitZwitscher)+"\n", res.Message)
}

func testNotEmptyZwitscher(t *testing.T) {
	res := ZwitscherListViewStruct{}
	acc := database.Account{}
	err := acc.GetByUserName("user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user"})
	res.validateZwitscherCreate(ctx, &database.Account{DisplayName: "user", ID: 1})

	assert.Equal(t, generics.ZwitscherIsNotAllowedToBeEmpty+"\n", res.Message)
}

func testThatAccountIsNotAllowedToWriteZwitscher(t *testing.T) {
	res := ZwitscherListViewStruct{}
	acc := database.Account{}
	err := acc.GetByUserName("user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "user"})
	res.validateZwitscherCreate(ctx, &database.Account{DisplayName: "lol", ID: 12})

	assert.Equal(t, generics.AccountIsNotYours+"\n", res.Message)

}

func testWriteAccountDoesNotExistInDatabaseZwitscher(t *testing.T) {
	res := ZwitscherListViewStruct{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	res.validateZwitscherCreate(ctx, &database.Account{})

	assert.Equal(t, generics.AccountDoesNotExists+"\n", res.Message)
}
