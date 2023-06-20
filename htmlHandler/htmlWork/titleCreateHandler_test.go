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
	"html/template"
	"testing"
)

func TestTitleCreateHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupTitleCreatePage", testSetupTitleCreatePage)
	t.Run("testGetTitleCreatePage", testGetTitleCreatePage)
	t.Run("testPostTitleCreatePage", testPostTitleCreatePage)
}

func testPostTitleCreatePage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testCreateTitleError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostCreateTitlePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testCreateTitle")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	PostCreateTitlePage(ctx)
	assert.Equal(t, "TestTitleCreate|{     []}|"+string(generics.NoMainGroupSubGroupOrNameProvided)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "createTitleExample", "mainGroup": "createTitleExample", "subGroup": "createTitleExample"})
	PostCreateTitlePage(ctx)
	assert.Equal(t, "TestTitleCreate|{ createTitleExample  createTitleExample createTitleExample []}|"+string(dataLogic.SuccessCreatedTitle)+"\n|true", w.Body.String())

	title := database.Title{}
	err = title.GetByName("createTitleExample")
	assert.Nil(t, err)
	assert.Equal(t, "createTitleExample", title.Name)
	assert.Equal(t, "createTitleExample", title.MainGroup)
	assert.Equal(t, "createTitleExample", title.SubGroup)
	assert.Equal(t, sql.NullString{
		String: "",
		Valid:  false,
	}, title.Flair)
	assert.Equal(t, 0, len(title.Holder))
}

func testGetTitleCreatePage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testCreateTitleError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetCreateTitlePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testCreateTitle")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetCreateTitlePage(ctx)
	assert.Equal(t, "TestTitleCreate|{     []}||false", w.Body.String())
}

func testSetupTitleCreatePage(t *testing.T) {
	Setup()
	htmlBasics.Setup()

	acc := database.Account{
		DisplayName: "testCreateTitleError",
		Username:    "testCreateTitleError",
		Role:        database.MediaAdmin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testCreateTitle",
		Username:    "testCreateTitle",
		Role:        database.Admin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestTitleCreate|{{.Page.Title}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "createTitle")
}
