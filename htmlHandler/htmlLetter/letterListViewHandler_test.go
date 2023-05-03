package htmlLetter

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
	"time"
)

func TestViewListOfLetters(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestLettersDB()
	database.TestAccountDB()

	t.Run("setupLettersForPaging", setupLettersForPaging)
	t.Run("setupPagesForTestViewListOfLetters", setupPagesForTestViewListOfLetters)
	t.Run("testViewModMailPage", testViewModMailPage)
	t.Run("testPostViewLetterPage", testPostViewLetterPage)
	t.Run("testViewLetterPage", testViewLetterPage)
	t.Run("testNotAuthorized", testNotAuthorized)
	t.Run("testAccountNotYours", testAccountNotYours)
}

func testAccountNotYours(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=askdbald"
	GetViewLetterListPage(ctx)
	assert.Equal(t, "error "+generics.AccountDoesNotExistOrIsNotYours, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "usr=admin"
	GetViewLetterListPage(ctx)
	assert.Equal(t, "error "+generics.AccountDoesNotExistOrIsNotYours, w.Body.String())
}

func testNotAuthorized(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostViewLetterListPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	w, ctx = htmlHandler.GetEmptyContext(t)
	GetViewLetterListPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	w, ctx = htmlHandler.GetEmptyContext(t)
	GetViewModMailListPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
}

func testViewLetterPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test5&amount=12"
	GetViewLetterListPage(ctx)
	assert.Equal(t, "letterList true false true  test4 12  test4 d true test3 c true test2 b true test a true", w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test3&type=before&amount=1"
	GetViewLetterListPage(ctx)
	assert.Equal(t, "letterList true true true test4 test4 1  test4 d true", w.Body.String())
}

func testPostViewLetterPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"selectedAccount": "bazinga"})
	PostViewLetterListPage(ctx)
	assert.Equal(t, "/letter-list?usr=bazinga", w.Header().Get("Location"))
}

func testViewModMailPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test5&amount=12"
	GetViewModMailListPage(ctx)
	assert.Equal(t, "letterList false false true  test4 12  test4 d true test3 c true test2 b true test a true", w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test3&type=before&amount=1"
	GetViewModMailListPage(ctx)
	assert.Equal(t, "letterList false true true test4 test4 1  test4 d true", w.Body.String())
}

func setupPagesForTestViewListOfLetters(t *testing.T) {
	acc := database.Account{
		ID:            0,
		DisplayName:   "test",
		Flair:         "",
		Username:      "test",
		Password:      "",
		Suspended:     false,
		RefToken:      sql.NullString{},
		ExpDate:       sql.NullTime{},
		LoginTries:    0,
		NextLoginTime: sql.NullTime{},
		Role:          database.User,
		Linked:        sql.NullInt64{},
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username, acc.Role = "admin", "admin", database.MediaAdmin
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Search}} {{.Page.HasNext}} {{.Page.HasBefore}} {{.Page.NextUUID}} {{.Page.BeforeUUID}} {{.Page.Amount}} {{range $i, $pub := .Page.LetterList}} {{$pub.UUID}} {{$pub.Title}} {{$pub.Info.ModMessage}}{{end}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":      temp,
			"letterList": temp2,
		},
	}
}

func TestValidateLetterPaging(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestLettersDB()

	t.Run("testEmptyPages", testEmptyPages)
	t.Run("setupLettersForPaging", setupLettersForPaging)
	t.Run("testValidateLetterReadNextPage", testValidateLetterReadNextPage)
	t.Run("testValidateLetterReadPageBefore", testValidateLetterReadPageBefore)
}

func testValidateLetterReadPageBefore(t *testing.T) {
	first := ViewLetterListStruct{}
	second := ViewLetterListStruct{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test3"
	err := first.validateLetterReadPageBefore(ctx, 1, "", true)
	assert.Nil(t, err)
	err = second.validateLetterReadPageBefore(ctx, 1, "test", false)
	assert.Nil(t, err)
	assert.Equal(t, first, second)
	assert.True(t, first.HasNext)
	assert.True(t, first.HasBefore)
	assert.Equal(t, 1, len(first.LetterList))
	assert.Equal(t, "test4", first.NextUUID)
	assert.Equal(t, "test4", first.BeforeUUID)
}

func testValidateLetterReadNextPage(t *testing.T) {
	first := ViewLetterListStruct{}
	second := ViewLetterListStruct{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test3"
	err := first.validateLetterReadNextPage(ctx, 1, "", true)
	assert.Nil(t, err)
	err = second.validateLetterReadNextPage(ctx, 1, "test", false)
	assert.Nil(t, err)
	assert.Equal(t, first, second)
	assert.True(t, first.HasNext)
	assert.True(t, first.HasBefore)
	assert.Equal(t, 1, len(first.LetterList))
	assert.Equal(t, "test2", first.NextUUID)
	assert.Equal(t, "test2", first.BeforeUUID)
}

func setupLettersForPaging(t *testing.T) {
	letter := database.Letter{
		UUID:        "test",
		Author:      "Bazinga",
		Flair:       "a",
		Title:       "a",
		Content:     "a",
		HTMLContent: "a",
		Info: database.LetterInfo{
			PeopleInvitedToSign: []string{"test"},
			ModMessage:          true,
		},
	}
	err := letter.CreateMe()
	assert.Nil(t, err)
	time.Sleep(100)
	letter.UUID, letter.Title = "test2", "b"
	err = letter.CreateMe()
	assert.Nil(t, err)
	time.Sleep(100)
	letter.UUID, letter.Title = "test3", "c"
	err = letter.CreateMe()
	assert.Nil(t, err)
	time.Sleep(100)
	letter.UUID, letter.Title = "test4", "d"
	err = letter.CreateMe()
	assert.Nil(t, err)
	time.Sleep(100)
	letter.UUID, letter.Title = "test5", "e"
	err = letter.CreateMe()
	assert.Nil(t, err)
}

func testEmptyPages(t *testing.T) {
	listStruct := ViewLetterListStruct{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	err := listStruct.validateLetterReadNextPage(ctx, 10, "", true)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(listStruct.LetterList))
	assert.Equal(t, false, listStruct.HasNext)
	assert.Equal(t, false, listStruct.HasBefore)
}
