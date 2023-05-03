package htmlWork

import (
	"API_MBundestag/dataLogic"
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
)

func TestPagesTitleEdit(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestTitlesDB()

	t.Run("setupFillEditTitleStruct", setupFillEditTitleStruct)
	t.Run("setupPagesTitleEdit", setupPagesTitleEdit)
	t.Run("testGetEditTitlePage", testGetEditTitlePage)
	t.Run("PostEditTitlePageSearch", PostEditTitlePageSearch)
	t.Run("testPostEditTitlePage", testPostEditTitlePage)
	t.Run("PostEditTitlePageDelete", PostEditTitlePageDelete)
}

func PostEditTitlePageDelete(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testChange"})
	ctx.Request.URL.RawQuery = "type=delete"
	PostEditTitlePage(ctx)
	assert.Equal(t, "editTitle [] [testChange] [testChange] testChange testChange testChange [admin]"+generics.SuccesfulDeletedTitle+"\n", w.Body.String())
	title := database.Title{}
	err = title.GetByName("changeTest")
	assert.Equal(t, sql.ErrNoRows, err)
}

func testPostEditTitlePage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "test", "newName": "testChange", "mainGroup": "testChange", "subGroup": "testChange", "user": "admin"})
	PostEditTitlePage(ctx)
	assert.Equal(t, "editTitle [test] [test testChange] [test testChange] testChange testChange testChange [admin]"+generics.SuccessFullEditTitle+"\n", w.Body.String())
	title := database.Title{}
	err = title.GetByName("testChange")
	assert.Nil(t, err)
	assert.Equal(t, "testChange", title.MainGroup)
	assert.Equal(t, "testChange", title.SubGroup)
	assert.Equal(t, "admin", title.Info.Names[0])
}

func PostEditTitlePageSearch(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostEditTitlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "test"})
	ctx.Request.URL.RawQuery = "type=search"
	PostEditTitlePage(ctx)
	assert.Equal(t, "editTitle [test] [test] [test] test test test []"+generics.SuccessFullFoundTitle+"\n", w.Body.String())
}

func testGetEditTitlePage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetEditTitlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetEditTitlePage(ctx)
	assert.Equal(t, "editTitle [test] [test] [test]    []", w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "title=test"
	GetEditTitlePage(ctx)
	assert.Equal(t, "editTitle [test] [test] [test] test test test []", w.Body.String())
}

func setupPagesTitleEdit(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.TitleNames}} {{.Page.ExistingMainGroup}} {{.Page.ExistingSubGroup}} {{.Page.Title.Name}} {{.Page.Title.MainGroup}} {{.Page.Title.SubGroup}} {{.Page.Title.Info.Names}}{{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":     temp,
			"editTitle": temp2,
		},
	}
}

func TestValidateTitelFuncs(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestTitlesDB()

	t.Run("setupFillEditTitleStruct", setupFillEditTitleStruct)
	t.Run("testTitleDoesNotExists", testTitleDoesNotExists)
	t.Run("testNoMainGroupSubGroupOrNameProvidedEditTitle", testNoMainGroupSubGroupOrNameProvidedEditTitle)
	t.Run("testUserAccountDoesNotExistErrorEditTitle", testUserAccountDoesNotExistErrorEditTitle)
	t.Run("testSuccessFullEditTitle", testSuccessFullEditTitle)

	t.Run("testTitleDoesNotExistsSearchTitle", testTitleDoesNotExistsSearchTitle)
	t.Run("testFindTitleSearch", testFindTitleSearch)

	t.Run("testTitleDoesNotExistsDeleteTitle", testTitleDoesNotExistsDeleteTitle)
	t.Run("testDeleteTitle", testDeleteTitle)
}

func testDeleteTitle(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "testChange"})
	title := database.Title{}
	err := title.GetByName("testChange")
	assert.Nil(t, err)
	res := validateDeleteTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title:             title,
		ExistingMainGroup: []string{"testChange"},
		ExistingSubGroup:  []string{"testChange"},
		TitleNames:        []string{},
		Names:             []string{"admin"},
		Message:           generics.SuccesfulDeletedTitle + "\n",
	}, *res)
	assert.Equal(t, []string{}, dataLogic.MainGroupNames)
	assert.Equal(t, []string{}, dataLogic.SubGroupNames)
	assert.Equal(t, []string{}, dataLogic.TitleNames)
	err = title.GetByName("testChange")
	assert.Equal(t, sql.ErrNoRows, err)
}

func testTitleDoesNotExistsDeleteTitle(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := validateDeleteTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title:             database.Title{},
		ExistingMainGroup: []string{"testChange"},
		ExistingSubGroup:  []string{"testChange"},
		TitleNames:        []string{"testChange"},
		Names:             []string{"admin"},
		Message:           generics.TitleDoesNotExists + "\n",
	}, *res)
}

func testFindTitleSearch(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "testChange"})
	res := validateSearchTitle(ctx)
	title := database.Title{}
	err := title.GetByName("testChange")
	assert.Nil(t, err)
	assert.Equal(t, EditTitleStruct{
		Title:             title,
		ExistingMainGroup: []string{"testChange"},
		ExistingSubGroup:  []string{"testChange"},
		TitleNames:        []string{"testChange"},
		Names:             []string{"admin"},
		Message:           generics.SuccessFullFoundTitle + "\n",
	}, *res)
}

func testTitleDoesNotExistsSearchTitle(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := validateSearchTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title:             database.Title{},
		ExistingMainGroup: []string{"testChange"},
		ExistingSubGroup:  []string{"testChange"},
		TitleNames:        []string{"testChange"},
		Names:             []string{"admin"},
		Message:           generics.TitleDoesNotExists + "\n",
	}, *res)
}

func testSuccessFullEditTitle(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "test", "newName": "testChange", "mainGroup": "testChange", "subGroup": "testChange", "user": "admin"})
	res := validateEditTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title: database.Title{
			Name:      "testChange",
			MainGroup: "testChange",
			SubGroup:  "testChange",
			Info:      database.TitleInfo{Names: []string{"admin"}},
		},
		ExistingMainGroup: []string{"test", "testChange"},
		ExistingSubGroup:  []string{"test", "testChange"},
		TitleNames:        []string{"test"},
		Names:             []string{"admin"},
		Message:           generics.SuccessFullEditTitle + "\n",
	}, *res)
}

func testUserAccountDoesNotExistErrorEditTitle(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "test", "newName": "testChange", "mainGroup": "testChange", "subGroup": "testChange", "user": "lol"})
	res := validateEditTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title: database.Title{
			Name:      "test",
			MainGroup: "testChange",
			SubGroup:  "testChange",
			Info:      database.TitleInfo{Names: []string{"lol"}},
		},
		ExistingMainGroup: []string{"test"},
		ExistingSubGroup:  []string{"test"},
		TitleNames:        []string{"test"},
		Names:             []string{"admin"},
		Message:           fmt.Sprintf(generics.AccountDoesNotExistError, "lol") + "\n",
	}, *res)
}

func testNoMainGroupSubGroupOrNameProvidedEditTitle(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "test"})
	res := validateEditTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title: database.Title{
			Name: "test",
			Info: database.TitleInfo{Names: []string{}},
		},
		ExistingMainGroup: []string{"test"},
		ExistingSubGroup:  []string{"test"},
		TitleNames:        []string{"test"},
		Names:             []string{"admin"},
		Message:           generics.NoMainGroupSubGroupOrNameProvided + "\n",
	}, *res)
}

func testTitleDoesNotExists(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := validateEditTitle(ctx)
	assert.Equal(t, EditTitleStruct{
		Title:             database.Title{Info: database.TitleInfo{Names: []string{}}},
		ExistingMainGroup: []string{"test"},
		ExistingSubGroup:  []string{"test"},
		TitleNames:        []string{"test"},
		Names:             []string{"admin"},
		Message:           generics.TitleDoesNotExists + "\n",
	}, *res)
}

func TestGetEmptyEditTitleStruct(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	dataLogic.MainGroupNames = nil
	dataLogic.SubGroupNames = nil
	dataLogic.TitleNames = nil

	database.TestAccountDB()
	database.TestTitlesDB()

	t.Run("setupFillEditTitleStruct", setupFillEditTitleStruct)
	t.Run("testCorrectlyFillEditTitleStruct", testCorrectlyFillEditTitleStruct)
}

func testCorrectlyFillEditTitleStruct(t *testing.T) {
	res := getEmptyEditTitleStruct()
	assert.Equal(t, EditTitleStruct{
		ExistingMainGroup: []string{"test"},
		ExistingSubGroup:  []string{"test"},
		TitleNames:        []string{"test"},
		Names:             []string{"admin"},
		Message:           "",
	}, *res)
}

func setupFillEditTitleStruct(t *testing.T) {
	htmlHandler.CreateAccountForTest(t, "admin", database.Admin, 0)
	title := database.Title{
		Name:      "test",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{},
		Info:      database.TitleInfo{Names: []string{}},
	}
	err := title.CreateMe()
	assert.Nil(t, err)
	err = dataLogic.RefreshTitleHierarchy()
	assert.Nil(t, err)
}
