package htmlLetter

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
	"strings"
	"testing"
	"time"
)

func TestLetterCreateHTMLFunctions(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestLettersDB()

	t.Run("setupAccountsAndLetters", setupAccountsAndLetters)
	t.Run("setupTemplatesForLetterCreate", setupTemplatesForLetterCreate)
	t.Run("testGetLetterCreate", testGetLetterCreate)
	t.Run("testGetModMailCreate", testGetModMailCreate)
	t.Run("testPostLetterCreate", testPostLetterCreate)
	t.Run("testPostModMailCreate", testPostModMailCreate)
}

func testPostModMailCreate(t *testing.T) {
	//fail
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	PostCreateModMailPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
	//test preview
	acc = database.Account{}
	err = acc.GetByUserName("test_admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "# Headline"})
	ctx.Request.URL.RawQuery = "type=preview"
	PostCreateModMailPage(ctx)
	assert.Equal(t, "createLetter  true false <h1 "+helper.ReplacerMap["h1"]+">Headline</h1>\n"+generics.PreviewText+"\n", w.Body.String())
	//test fail send
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "# Headline"})
	PostCreateModMailPage(ctx)
	assert.Equal(t, "createLetter  true false "+generics.AuthorEmptyError+"\n", w.Body.String())
	//test successfull
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "test", "title": "test", "author": "Test Author"})
	PostCreateModMailPage(ctx)
	truth := strings.HasPrefix(w.Result().Header.Get("Location"), "/letter?uuid=") &&
		strings.HasSuffix(w.Result().Header.Get("Location"), "&usr=test_admin")
	assert.True(t, truth)
}

func testPostLetterCreate(t *testing.T) {
	//fail
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostCreateLetterPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
	//test preview
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "# Headline"})
	ctx.Request.URL.RawQuery = "type=preview"
	PostCreateLetterPage(ctx)
	assert.Equal(t, "createLetter  false false <h1 "+helper.ReplacerMap["h1"]+">Headline</h1>\n"+generics.PreviewText+"\n", w.Body.String())
	//test fail send
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "# Headline"})
	PostCreateLetterPage(ctx)
	assert.Equal(t, "createLetter  false false "+generics.ContentAndTitelAreEmpty+"\n", w.Body.String())
	//test successfull
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "test", "title": "test", "selectedAccount": "test_user"})
	PostCreateLetterPage(ctx)
	truth := strings.HasPrefix(w.Result().Header.Get("Location"), "/letter?uuid=") &&
		strings.HasSuffix(w.Result().Header.Get("Location"), "&usr=test_user")
	assert.True(t, truth)
}

func testGetModMailCreate(t *testing.T) {
	//fail
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	GetCreateModMailPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
	//success
	err = acc.GetByUserName("test_admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetCreateModMailPage(ctx)
	assert.Equal(t, "createLetter  true true ", w.Body.String())
}

func testGetLetterCreate(t *testing.T) {
	//fail
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetCreateLetterPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
	//success
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetCreateLetterPage(ctx)
	assert.Equal(t, "createLetter  false true ", w.Body.String())
}

func setupTemplatesForLetterCreate(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Letter.Title}} {{.Page.ModMail}} {{.Page.Letter.Info.NoSigning}} {{.Page.Preview}}{{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":        temp,
			"createLetter": temp2,
		},
	}
	helper.UpdateAttributes()
}

func TestLetterValidator(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestLettersDB()

	t.Run("setupAccountsAndLetters", setupAccountsAndLetters)
	t.Run("testEmptyAuthor", testEmptyAuthor)
	t.Run("testEmptyTextOrTitle", testEmptyTextOrTitle)
	t.Run("testAccountDoesNotExist", testAccountDoesNotExist)
	t.Run("testAccountNotAllowedToPostLetter", testAccountNotAllowedToPostLetter)
	t.Run("testAccountListedDoesNotExist", testAccountListedDoesNotExist)
	t.Run("testDifferentCreates", testDifferentCreates)
}

func testDifferentCreates(t *testing.T) {
	//test with user not in the list
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"content": "test", "title": "test", "selectedAccount": "test_user", "user": "test_admin"})
	letterStruct, err := validateLetterCreate(ctx, &acc, false)
	assert.Nil(t, err)
	assert.Equal(t, "", letterStruct.Message)
	letter := database.Letter{}
	err = letter.GetByID(letterStruct.Letter.UUID)
	letter.Written = time.Time{}
	assert.Equal(t, letter, letterStruct.Letter)
	assert.Equal(t, 2, len(letter.Info.PeopleInvitedToSign))
	assert.Equal(t, "test_admin", letter.Info.PeopleInvitedToSign[0])
	assert.Equal(t, "test_user", letter.Info.PeopleInvitedToSign[1])
	assert.Equal(t, 1, len(letter.Info.PeopleNotYetSigned))
	assert.Equal(t, "test_admin", letter.Info.PeopleNotYetSigned[0])
	assert.Equal(t, 1, len(letter.Info.Signed))
	assert.Equal(t, "test_user", letter.Info.Signed[0])
	//test with user in the list
	err = acc.GetByUserName("test_user")
	assert.Nil(t, err)
	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"content": "test", "title": "test", "selectedAccount": "test_user", "user": "test_user"})
	ctx.Request.PostForm.Add("user", "test_admin")
	letterStruct, err = validateLetterCreate(ctx, &acc, false)
	assert.Nil(t, err)
	assert.Equal(t, "", letterStruct.Message)
	err = letter.GetByID(letterStruct.Letter.UUID)
	letter.Written = time.Time{}
	assert.Equal(t, letter, letterStruct.Letter)
	assert.Equal(t, 2, len(letter.Info.PeopleInvitedToSign))
	assert.Equal(t, "test_user", letter.Info.PeopleInvitedToSign[0])
	assert.Equal(t, "test_admin", letter.Info.PeopleInvitedToSign[1])
	assert.Equal(t, 1, len(letter.Info.PeopleNotYetSigned))
	assert.Equal(t, "test_admin", letter.Info.PeopleNotYetSigned[0])
	assert.Equal(t, 1, len(letter.Info.Signed))
	assert.Equal(t, "test_user", letter.Info.Signed[0])
	//test no signing
	err = acc.GetByUserName("test_user")
	assert.Nil(t, err)
	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"content": "test", "title": "test", "selectedAccount": "test_user", "noSigning": "true"})
	letterStruct, err = validateLetterCreate(ctx, &acc, false)
	assert.Nil(t, err)
	assert.Equal(t, "", letterStruct.Message)
	err = letter.GetByID(letterStruct.Letter.UUID)
	letter.Written = time.Time{}
	assert.Equal(t, letter, letterStruct.Letter)
	assert.Equal(t, 0, len(letter.Info.Signed))
	assert.Equal(t, 0, len(letter.Info.PeopleNotYetSigned))
	assert.Equal(t, 0, len(letter.Info.Rejected))
	//test signing with modmail
	err = acc.GetByUserName("test_admin")
	assert.Nil(t, err)
	_, ctx = htmlHandler.GetContextWithForm(t, map[string]string{"author": "test", "content": "test", "title": "test", "selectedAccount": "test_admin", "user": "test_user"})
	letterStruct, err = validateLetterCreate(ctx, &acc, true)
	assert.Nil(t, err)
	assert.Equal(t, "", letterStruct.Message)
	err = letter.GetByID(letterStruct.Letter.UUID)
	letter.Written = time.Time{}
	assert.Equal(t, letter, letterStruct.Letter)
	assert.Equal(t, 0, len(letter.Info.Signed))
	assert.Equal(t, 1, len(letter.Info.PeopleNotYetSigned))
	assert.Equal(t, "test_user", letter.Info.PeopleNotYetSigned[0])
	assert.Equal(t, 1, len(letter.Info.PeopleInvitedToSign))
	assert.Equal(t, "test_user", letter.Info.PeopleInvitedToSign[0])
}

func testAccountListedDoesNotExist(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"content": "test", "title": "test", "selectedAccount": "test_user", "user": "bazinga"})
	letterStruct, err := validateLetterCreate(ctx, &acc, false)
	assert.Equal(t, ErrorInLetter, err)
	assert.Equal(t, fmt.Sprintf(generics.AccountDoesNotExistError, "bazinga")+"\n", letterStruct.Message)
}

func testAccountNotAllowedToPostLetter(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"content": "test", "title": "test", "selectedAccount": "test_user2"})
	letterStruct, err := validateLetterCreate(ctx, &acc, false)
	assert.Equal(t, ErrorInLetter, err)
	assert.Equal(t, generics.AccountIsNotYours+"\n", letterStruct.Message)
}

func testAccountDoesNotExist(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"content": "test", "title": "test"})
	letterStruct, err := validateLetterCreate(ctx, &database.Account{}, false)
	assert.Equal(t, ErrorInLetter, err)
	assert.Equal(t, generics.AccountDoesNotExists+"\n", letterStruct.Message)
}

func testEmptyTextOrTitle(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	letterStruct, err := validateLetterCreate(ctx, &database.Account{}, false)
	assert.Equal(t, ErrorInLetter, err)
	assert.Equal(t, generics.ContentAndTitelAreEmpty+"\n", letterStruct.Message)
}

func testEmptyAuthor(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	letterStruct, err := validateLetterCreate(ctx, &database.Account{}, true)
	assert.Equal(t, ErrorInLetter, err)
	assert.Equal(t, generics.AuthorEmptyError+"\n", letterStruct.Message)
}

func setupAccountsAndLetters(t *testing.T) {
	acc := database.Account{
		DisplayName:   "test_user",
		Flair:         "",
		Username:      "test_user",
		Password:      "test_user",
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
	acc.DisplayName, acc.Username = "test_user2", "test_user2"
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "test_admin", "test_admin"
	acc.Role = database.Admin
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "test_press", "test_press"
	acc.Role = database.PressAccount
	acc.Linked.Valid = true
	acc.Linked.Int64 = 1
	err = acc.CreateMe()
	assert.Nil(t, err)
}

func TestLetterStructGenerator(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestLettersDB()

	t.Run("setupAccountsAndLetters", setupAccountsAndLetters)
	t.Run("testGenerateStructLetterCreate", testGenerateStructLetterCreate)
	t.Run("testGenerateStructLetterCreateWithMessage", testGenerateStructLetterCreateWithMessage)
}

func testGenerateStructLetterCreate(t *testing.T) {
	letter := getEmtpyLetterCreateStruct(false, &database.Account{})
	assert.Equal(t, LetterCreatePageStruct{
		Message:         "",
		Names:           []string{"test_admin", "test_press", "test_user", "test_user2"},
		SelectedAccount: "",
		Accounts:        []database.Account{{}},
		Preview:         "",
		Letter:          database.Letter{Info: database.LetterInfo{NoSigning: true}},
		ModMail:         false,
	}, *letter)
	letter = getEmtpyLetterCreateStruct(true, &database.Account{ID: 1})
	acc := database.Account{}
	err := acc.GetByUserName("test_press")
	assert.Nil(t, err)
	assert.Equal(t, LetterCreatePageStruct{
		Message:         "",
		Names:           []string{"test_admin", "test_press", "test_user", "test_user2"},
		SelectedAccount: "",
		Accounts:        []database.Account{{ID: 1}, acc},
		Preview:         "",
		Letter:          database.Letter{Info: database.LetterInfo{NoSigning: true}},
		ModMail:         true,
	}, *letter)
}

func testGenerateStructLetterCreateWithMessage(t *testing.T) {
	letter := getEmtpyLetterCreateStructWithMessage(false, &database.Account{}, "testsdf")
	assert.Equal(t, LetterCreatePageStruct{
		Message:         "testsdf\n",
		Names:           []string{"test_admin", "test_press", "test_user", "test_user2"},
		SelectedAccount: "",
		Accounts:        []database.Account{{}},
		Preview:         "",
		Letter:          database.Letter{Info: database.LetterInfo{NoSigning: true}},
		ModMail:         false,
	}, *letter)
	letter = getEmtpyLetterCreateStructWithMessage(true, &database.Account{ID: 1}, "a")
	acc := database.Account{}
	err := acc.GetByUserName("test_press")
	assert.Nil(t, err)
	assert.Equal(t, LetterCreatePageStruct{
		Message:         "a\n",
		Names:           []string{"test_admin", "test_press", "test_user", "test_user2"},
		SelectedAccount: "",
		Accounts:        []database.Account{{ID: 1}, acc},
		Preview:         "",
		Letter:          database.Letter{Info: database.LetterInfo{NoSigning: true}},
		ModMail:         true,
	}, *letter)
}
