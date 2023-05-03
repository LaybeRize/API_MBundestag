package htmlWork

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
)

func TestDatabaseFailureOrgView(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestDatabase("DROP TABLE IF EXISTS organisations;", "DROP TYPE IF EXISTS STATUS;")

	t.Run("testFailingDatabaseOrgView", testFailingDatabaseOrgView)
}

func testFailingDatabaseOrgView(t *testing.T) {
	acc := database.Account{
		DisplayName:   "test_user",
		Flair:         "",
		Username:      "test_user",
		Password:      "asfcv",
		Suspended:     false,
		RefToken:      sql.NullString{},
		ExpDate:       sql.NullTime{},
		LoginTries:    0,
		NextLoginTime: sql.NullTime{},
		Role:          database.Admin,
		Linked:        sql.NullInt64{},
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error": temp2,
		},
	}

	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	GetHiddenOrganisationViewPage(ctx)
	assert.Equal(t, "error "+generics.ErrorWhileLodingOrganisationView, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetOrganisationViewPage(ctx)
	assert.Equal(t, "error "+generics.ErrorWhileLodingOrganisationView, w.Body.String())
}

func TestOrgansationView(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupTestOrgView", setupTestOrgView)
	t.Run("testWithoutLogin", testWithoutLogin)
	t.Run("testWithAccount", testWithAccount)
	t.Run("testFailAdminView", testFailAdminView)
	t.Run("testAdminView", testAdminView)
}

func testAdminView(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	GetHiddenOrganisationViewPage(ctx)
	assert.Equal(t, "hiddenOrganisation asd 2 a 1 sasdad bnfasd 1 fwe ", w.Body.String())
}

func testFailAdminView(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	GetHiddenOrganisationViewPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())
}

func testWithAccount(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test_user")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUser(t, acc)
	GetOrganisationViewPage(ctx)
	assert.Equal(t, "displayOrganisation a 0 a 1 a asd 0 a 2 bazinga test  test_user", w.Body.String())
}

func testWithoutLogin(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetOrganisationViewPage(ctx)
	assert.Equal(t, "displayOrganisation a 0 a 1 a asd 0 a 1 bazinga ", w.Body.String())
}

func setupTestOrgView(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{range $i, $main := .Page}}{{$main.Name}} {{$main.Amount}} {{range $j, $sub := $main.Groups}}{{$sub.Name}} {{$sub.Amount}} {{range $k, $title := $sub.Organisations}}{{$title.Name}} {{range $z, $name := $title.Info.Viewer}} {{$name}}{{end}}{{end}}{{end}}{{end}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"displayOrganisation": temp,
			"hiddenOrganisation":  temp,
			"error":               temp2,
		},
	}
	acc := database.Account{
		DisplayName:   "test_user",
		Flair:         "",
		Username:      "test_user",
		Password:      "asfcv",
		Suspended:     false,
		RefToken:      sql.NullString{},
		ExpDate:       sql.NullTime{},
		LoginTries:    0,
		NextLoginTime: sql.NullTime{},
		Role:          database.User,
		Linked:        sql.NullInt64{},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username, acc.Role = "test_admin", "test_admin", database.Admin
	err = acc.CreateMe()
	assert.Nil(t, err)
	org := database.Organisation{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Flair:     sql.NullString{},
		Status:    database.Public,
		Info:      database.OrganisationInfo{},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org.MainGroup, org.Name = "asd", "bazinga"
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.Status = "test", database.Secret
	org.Info.Viewer = []string{"test_user"}
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Info.Viewer = []string{}
	org.Status = database.Hidden
	org.Name = "sasdad"
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name = "fwe"
	org.SubGroup = "bnfasd"
	err = org.CreateMe()
}
