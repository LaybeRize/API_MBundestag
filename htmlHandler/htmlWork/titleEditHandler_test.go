package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"html/template"
	"testing"
)

func TestTitleEditHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupTitleEditPage", testSetupTitleEditPage)
	t.Run("testGetTitleEditPage", testGetTitleEditPage)
	t.Run("testPostTitleEditPageSearch", testPostTitleEditPageSearch)
	t.Run("testPostTitleEditPage", testPostTitleEditPage)
	t.Run("testPostTitleEditPageDelete", testPostTitleEditPageDelete)
}

func testPostTitleEditPageDelete(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditTitle")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "sdkhbaldhsad"}, "type=delete")
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{sdkhbaldhsad     []}|"+string(dataLogic.ErrorCouldNotFindTitle)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "testEditHTML2"}, "type=delete")
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{testEditHTML2     []}|"+string(dataLogic.SuccessDeletedTitle)+"\n|true", w.Body.String())

	title := database.Title{}
	err = title.GetByName("testEditHTML")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = title.GetByName("testEditHTML2")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func testPostTitleEditPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditTitle")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUserWithForm(t, acc, map[string]string{})
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{     []}|"+string(generics.NoNameForTitleProvided)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testEditHTML", "newName": "testEditHTML2"})
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{testEditHTML testEditHTML2    []}|"+string(generics.NoMainGroupSubGroupOrNameProvided)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testEditHTML", "newName": "testEditHTML2", "subGroup": "testEditHTML2", "mainGroup": "testEditHTML2"})
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{testEditHTML2 testEditHTML2  testEditHTML2 testEditHTML2 []}|"+string(dataLogic.SuccessChangedTitle)+"\n|true", w.Body.String())

	title := database.Title{}
	err = title.GetByName("testEditHTML")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = title.GetByName("testEditHTML2")
	assert.Nil(t, err)
	assert.Equal(t, "testEditHTML2", title.Name)
	assert.Equal(t, "testEditHTML2", title.SubGroup)
	assert.Equal(t, "testEditHTML2", title.MainGroup)
	assert.Equal(t, 0, len(title.Holder))
	assert.Equal(t, sql.NullString{String: "", Valid: false}, title.Flair)
}

func testPostTitleEditPageSearch(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditTitleError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testEditTitle")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "sdkhbaldhsad"}, "type=search")
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{     []}|"+string(generics.TitleDoesNotExists)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "testEditHTML"}, "type=search")
	PostEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{testEditHTML testEditHTML  testEditHTML testEditHTML []}|"+string(generics.SuccessfulFoundTitle)+"\n|true", w.Body.String())
}

func testGetTitleEditPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditTitleError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetEditTitlePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testEditTitle")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetEditTitlePage(ctx)
	assert.Equal(t, "TestTitleEdit|{     []}||false", w.Body.String())
}

func testSetupTitleEditPage(t *testing.T) {
	Setup()
	htmlBasics.Setup()

	acc := database.Account{
		DisplayName: "testEditTitleError",
		Username:    "testEditTitleError",
		Role:        database.MediaAdmin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testEditTitle",
		Username:    "testEditTitle",
		Role:        database.Admin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	title := database.Title{
		Name:      "testEditHTML",
		MainGroup: "testEditHTML",
		SubGroup:  "testEditHTML",
	}
	err = title.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestTitleEdit|{{.Page.Title}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "editTitle")
}
