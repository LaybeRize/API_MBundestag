package htmlAccount

import (
	"API_MBundestag/database"
	gen "API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"html"
	"html/template"
	"testing"
)

func TestAccountListHandler(t *testing.T) {
	database.TestSetup()

	t.Run("setupSitesAccountListHandler", setupSitesAccountListHandler)
	t.Run("testViewAccountListHandler", testViewAccountListHandler)
}

func testViewAccountListHandler(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("c_test_admin_AccViewList")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetAdminListUserPage(ctx)
	assert.Equal(t, "TestError|"+gen.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("a_test_admin_AccViewList")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetAdminListUserPage(ctx)
	assert.Equal(t, "TestAdminViewUser|, a_test_admin_AccViewList, b_test_admin_AccViewList, c_test_admin_AccViewList, head_admin||false", html.UnescapeString(w.Body.String()))
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "acc=a_test_admin_AccViewList")
	GetAdminListUserPage(ctx)
	assert.Equal(t, "TestAdminViewUser|, a_test_admin_AccViewList, b_test_admin_AccViewList||false", html.UnescapeString(w.Body.String()))
}

func setupSitesAccountListHandler(t *testing.T) {
	acc := database.Account{
		DisplayName: "a_test_admin_AccViewList",
		Username:    "a_test_admin_AccViewList",
		Password:    "test_admin",
		Role:        database.HeadAdmin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	i := acc.ID
	acc = database.Account{
		DisplayName: "b_test_admin_AccViewList",
		Username:    "b_test_admin_AccViewList",
		Password:    "test_admin",
		Role:        database.PressAccount,
		Linked:      sql.NullInt64{Valid: true, Int64: i},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "c_test_admin_AccViewList",
		Username:    "c_test_admin_AccViewList",
		Password:    "test_admin",
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(template.FuncMap{
		"testFunc": func(arr database.AccountList) string {
			str := ""
			for _, a := range arr {
				switch a.Username {
				case "a_test_admin_AccViewList", "b_test_admin_AccViewList", "c_test_admin_AccViewList", "d_test_admin_AccViewList", "head_admin":
					str += ", " + a.Username
				}
			}
			return str
		},
	}).Parse("TestAdminViewUser|{{testFunc .Page.Accounts}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "adminViewUser")
}
