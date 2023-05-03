package htmlAccount

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"testing"
)

func TestAccountSelfHandler(t *testing.T) {
	database.TestSetup()
	t.Run("setupSitesAccountSelfHandler", setupSitesAccountSelfHandler)
	t.Run("testSelfViewHandling", testSelfViewHandling)
	t.Run("testGetChangePasswordHandling", testGetChangePasswordHandling)
	t.Run("testPostChangePasswordHandling", testPostChangePasswordHandling)
}

func testPostChangePasswordHandling(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	PostPasswordChangePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	//check if new must be the same
	acc := database.Account{}
	err := acc.GetByUserName("a_test_admin_AccSelfView")
	oldPassword := acc.Password
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"newPassword": "newPassword"})
	PostPasswordChangePage(ctx)
	assert.Equal(t, "TestPasswordChange|"+string(generics.NewPasswordIsNotTheSame)+"|false", w.Body.String())
	err = acc.GetByUserName("a_test_admin_AccSelfView")
	assert.Nil(t, err)
	assert.Equal(t, oldPassword, acc.Password)
	//check min length of new password
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"newPassword": "a", "newPassword2": "a"})
	PostPasswordChangePage(ctx)
	assert.Equal(t, "TestPasswordChange|"+string(generics.NewPasswordIsNotMinimumOf10Characters)+"|false", w.Body.String())
	err = acc.GetByUserName("a_test_admin_AccSelfView")
	assert.Nil(t, err)
	assert.Equal(t, oldPassword, acc.Password)
	//check if old password is correct
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"newPassword": "1234567890", "newPassword2": "1234567890"})
	PostPasswordChangePage(ctx)
	assert.Equal(t, "TestPasswordChange|"+string(generics.OldPasswordNotcorrect)+"|false", w.Body.String())
	err = acc.GetByUserName("a_test_admin_AccSelfView")
	assert.Nil(t, err)
	assert.Equal(t, oldPassword, acc.Password)
	//success change password
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"password": "test", "newPassword": "1234567890", "newPassword2": "1234567890"})
	PostPasswordChangePage(ctx)
	assert.Equal(t, "TestPasswordChange|"+string(generics.SuccessChangePassword)+"|true", w.Body.String())
	err = acc.GetByUserName("a_test_admin_AccSelfView")
	assert.Nil(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte("1234567890"))
	assert.Nil(t, err)
}

func testGetChangePasswordHandling(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	GetPasswordChangePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("a_test_admin_AccSelfView")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetPasswordChangePage(ctx)
	assert.Equal(t, "TestPasswordChange||false", w.Body.String())
}

func testSelfViewHandling(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	GetViewOfProfilePage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("a_test_admin_AccSelfView")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetViewOfProfilePage(ctx)
	assert.Equal(t, "TestViewSelfUser|[{a_test_admin_AccSelfView testa  } {b_test_admin_AccSelfView testb test_AccSelfView }]", w.Body.String())
}

func setupSitesAccountSelfHandler(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "a_test_admin_AccSelfView",
		Username:    "a_test_admin_AccSelfView",
		Flair:       "testa",
		Password:    string(hash),
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	i := acc.ID
	acc = database.Account{
		DisplayName: "b_test_admin_AccSelfView",
		Username:    "b_test_admin_AccSelfView",
		Password:    "test_admin",
		Flair:       "testb",
		Role:        database.PressAccount,
		Linked:      sql.NullInt64{Valid: true, Int64: i},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	title := database.Title{
		Name:      "test_AccSelfView",
		MainGroup: "test_AccSelfView",
		SubGroup:  "test_AccSelfView",
		Flair:     sql.NullString{},
		Holder:    []database.Account{acc},
	}
	err = title.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestPasswordChange|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "password")
	temp, err = template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestViewSelfUser|{{.Page}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "viewPersonalInfo")
}
