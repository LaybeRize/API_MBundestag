package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestPagesTitleCreate(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestTitlesDB()
	err := dataLogic.RefreshTitleHierarchy()
	assert.Nil(t, err)

	t.Run("setupTitleCreatePage", setupTitleCreatePage)
	t.Run("setupAccountsTitleCreate", setupAccountsTitleCreate)

	t.Run("testGetTitleCreatePage", testGetTitleCreatePage)
	t.Run("testPostTitleCreatePage", testPostTitleCreatePage)
}

func testPostTitleCreatePage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostCreateTitlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "user": "test"})
	PostCreateTitlePage(ctx)
	assert.Equal(t, "createTitle [a] [a] a a a [test]"+generics.SuccessFullCreationTitle+"\n", w.Body.String())
}

func testGetTitleCreatePage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetCreateTitlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetCreateTitlePage(ctx)
	assert.Equal(t, "createTitle [] []    []", w.Body.String())
}

func setupTitleCreatePage(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.ExistingMainGroup}} {{.Page.ExistingSubGroup}} {{.Page.Title.Name}} {{.Page.Title.MainGroup}} {{.Page.Title.SubGroup}} {{.Page.Title.Info.Names}}{{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":       temp,
			"createTitle": temp2,
		},
	}
}

func TestValidateCreateTitle(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestTitlesDB()
	err := dataLogic.RefreshTitleHierarchy()
	assert.Nil(t, err)

	t.Run("setupAccountsTitleCreate", setupAccountsTitleCreate)
	t.Run("testNoMainGroupSubGroupOrNameProvidedTitle", testNoMainGroupSubGroupOrNameProvidedTitle)
	t.Run("testUserAccountDoesNotExistErrorTitle", testUserAccountDoesNotExistErrorTitle)
	t.Run("testSuccessFullCreationTitle", testSuccessFullCreationTitle)
}

func testSuccessFullCreationTitle(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "user": "test"})
	res := validateCreateTitle(ctx)
	assert.Equal(t, CreateTitleStruct{
		Title: database.Title{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info:      database.TitleInfo{Names: []string{"test"}},
		},
		ExistingMainGroup: []string{"a"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"admin", "press", "press2", "test", "test2"},
		Message:           generics.SuccessFullCreationTitle + "\n",
	}, *res)
	title := database.Title{}
	err := title.GetByName("a")
	assert.Nil(t, err)
	assert.Equal(t, database.Title{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Info:      database.TitleInfo{Names: []string{"test"}},
	}, title)
}

func testUserAccountDoesNotExistErrorTitle(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "user": "lol"})
	res := validateCreateTitle(ctx)
	assert.Equal(t, CreateTitleStruct{
		Title: database.Title{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info:      database.TitleInfo{Names: []string{"lol"}},
		},
		ExistingMainGroup: []string{},
		ExistingSubGroup:  []string{},
		Names:             []string{"admin", "press", "press2", "test", "test2"},
		Message:           fmt.Sprintf(generics.AccountDoesNotExistError, "lol") + "\n",
	}, *res)
}

func testNoMainGroupSubGroupOrNameProvidedTitle(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := validateCreateTitle(ctx)
	assert.Equal(t, CreateTitleStruct{
		Title:             database.Title{Info: database.TitleInfo{Names: []string{}}},
		ExistingMainGroup: []string{},
		ExistingSubGroup:  []string{},
		Names:             []string{"admin", "press", "press2", "test", "test2"},
		Message:           generics.NoMainGroupSubGroupOrNameProvided + "\n",
	}, *res)
}

func setupAccountsTitleCreate(t *testing.T) {
	htmlHandler.CreateAccountForTest(t, "admin", database.Admin, 0)
	htmlHandler.CreateAccountForTest(t, "test", database.User, 0)
	htmlHandler.CreateAccountForTest(t, "test2", database.User, 0)
	htmlHandler.CreateAccountForTest(t, "press", database.PressAccount, 2)
	htmlHandler.CreateAccountForTest(t, "press2", database.PressAccount, 3)
}

func TestGetEmptyCreateTitleStruct(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	dataLogic.MainGroupNames = nil
	dataLogic.SubGroupNames = nil

	database.TestAccountDB()
	dataLogic.MainGroupNames = []string{"test"}
	dataLogic.SubGroupNames = []string{"test"}

	t.Run("setupFillCreateTitleStruct", setupFillCreateTitleStruct)
	t.Run("testFailCorrectlyFillCreateTitleStruct", testCorrectlyFillCreateTitleStruct)
}

func testCorrectlyFillCreateTitleStruct(t *testing.T) {
	res := getEmptyCreateTitleStruct()
	assert.Equal(t, CreateTitleStruct{
		ExistingMainGroup: []string{"test"},
		ExistingSubGroup:  []string{"test"},
		Names:             []string{"admin"},
		Message:           "",
	}, *res)
}

func setupFillCreateTitleStruct(t *testing.T) {
	htmlHandler.CreateAccountForTest(t, "admin", database.Admin, 0)
}
