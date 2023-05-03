package htmlWork

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
	"net/http"
	"testing"
)

func TestSiteMakerOrCreate(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupCreateAccountsforOrgCreate", setupCreateAccountsforOrgCreate)
	t.Run("setupAdditionAccountsAndPageForOrgCreate", setupAdditionAccountsAndPageForOrgCreate)
	t.Run("testGetOrganisationCreate", testGetOrganisationCreate)
	t.Run("testPostCreateOrganisationPage", testPostCreateOrganisationPage)
}

func testPostCreateOrganisationPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostCreateOrganisationPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("test_admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"flair": "b", "name": "a", "mainGroup": "a", "subGroup": "a", "user": "test_press2", "admins": "test_press", "status": string(database.Public)})
	ctx.Request.PostForm.Add("admin", "test_press2")
	PostCreateOrganisationPage(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "createOrganisation [a] [a] a a a b {[test_press] [test_press2] [test_user test_user2]} "+string(database.Public)+generics.SuccessFullCreationOrg+"\n", w.Body.String())

	org := database.Organisation{}
	err = org.GetByName("a")
	assert.Nil(t, err)
	assert.Equal(t, database.Organisation{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Flair:     sql.NullString{Valid: true, String: "b"},
		Status:    database.Public,
		Info: database.OrganisationInfo{
			Admins: []string{"test_press"},
			User:   []string{"test_press2"},
			Viewer: []string{"test_user", "test_user2"},
		},
	}, org)
	err = acc.GetByDisplayName("test_press")
	assert.Nil(t, err)
	assert.Equal(t, "b", acc.Flair)
	err = acc.GetByDisplayName("test_press2")
	assert.Nil(t, err)
	assert.Equal(t, "b", acc.Flair)
}

func testGetOrganisationCreate(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetCreateOrganisationPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("test_admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetCreateOrganisationPage(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "createOrganisation [] []     {[] [] []} "+string(database.Public), w.Body.String())
}

func setupAdditionAccountsAndPageForOrgCreate(t *testing.T) {
	acc := database.Account{
		DisplayName: "test_admin",
		Username:    "test_admin",
		Password:    "test_admin",
		Role:        database.Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username, acc.Role = "test_press", "test_press", database.PressAccount
	acc.Linked = sql.NullInt64{Valid: true, Int64: 1}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username, acc.Linked.Int64 = "test_press2", "test_press2", 2
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.ExistingMainGroup}} {{.Page.ExistingSubGroup}} {{.Page.Organisation.Name}} {{.Page.Organisation.MainGroup}} {{.Page.Organisation.SubGroup}} {{.Page.Organisation.Flair.String}} {{.Page.Organisation.Info}} {{.Page.Organisation.Status}}{{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":              temp,
			"createOrganisation": temp2,
		},
	}
}

func TestValidateCreateOrgHandler(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupCreateAccountsforOrgCreate", setupCreateAccountsforOrgCreate)
	t.Run("testFailNoNames", testFailNoNames)
	t.Run("testFailAdminDoesNotExist", testFailAdminDoesNotExist)
	t.Run("testFailUserDoesNotExist", testFailUserDoesNotExist)
	t.Run("testFailStatusDoesNotExist", testFailStatusDoesNotExist)
	t.Run("testSuccessCreateBasic", testSuccessCreateBasic)
	t.Run("testSuccessCreateComplex", testSuccessCreateComplex)
}

func testSuccessCreateComplex(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"flair": "b", "name": "b", "mainGroup": "a", "subGroup": "a", "user": "test_user2", "admins": "test_user", "status": string(database.Public)})
	ctx.Request.PostForm.Add("admin", "test_user2")
	res := validateOrganisationCreate(ctx)
	assert.Equal(t, CreateOrganisationStruct{
		Message:           generics.SuccessFullCreationOrg + "\n",
		Names:             []string{"test_user", "test_user2"},
		ExistingMainGroup: []string{"a"},
		ExistingSubGroup:  []string{"a"},
		Organisation: database.Organisation{
			Name:      "b",
			MainGroup: "a",
			SubGroup:  "a",
			Flair:     sql.NullString{Valid: true, String: "b"},
			Status:    database.Public,
			Info: database.OrganisationInfo{
				Admins: []string{"test_user"},
				User:   []string{"test_user2"},
				Viewer: []string{"test_user", "test_user2"},
			},
		},
	}, *res)
	org := database.Organisation{}
	err := org.GetByName("b")
	assert.Nil(t, err)
	assert.Equal(t, database.Organisation{
		Name:      "b",
		MainGroup: "a",
		SubGroup:  "a",
		Flair:     sql.NullString{Valid: true, String: "b"},
		Status:    database.Public,
		Info: database.OrganisationInfo{
			Admins: []string{"test_user"},
			User:   []string{"test_user2"},
			Viewer: []string{"test_user", "test_user2"},
		},
	}, org)
	acc := database.Account{}
	err = acc.GetByDisplayName("test_user")
	assert.Nil(t, err)
	assert.Equal(t, "b", acc.Flair)
	err = acc.GetByDisplayName("test_user2")
	assert.Nil(t, err)
	assert.Equal(t, "b", acc.Flair)
}

func testSuccessCreateBasic(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "status": string(database.Secret)})
	res := validateOrganisationCreate(ctx)
	assert.Equal(t, CreateOrganisationStruct{
		Message:           generics.SuccessFullCreationOrg + "\n",
		Names:             []string{"test_user", "test_user2"},
		ExistingMainGroup: []string{"a"},
		ExistingSubGroup:  []string{"a"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Status:    database.Secret,
			Info: database.OrganisationInfo{
				Admins: []string{},
				User:   []string{},
				Viewer: []string{},
			},
		},
	}, *res)
	org := database.Organisation{}
	err := org.GetByName("a")
	assert.Nil(t, err)
	assert.Equal(t, database.Organisation{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Status:    database.Secret,
		Info: database.OrganisationInfo{
			Admins: []string{},
			User:   []string{},
			Viewer: []string{},
		},
	}, org)
}

func testFailStatusDoesNotExist(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a"})
	res := validateOrganisationCreate(ctx)
	assert.Equal(t, CreateOrganisationStruct{
		Message: generics.StatusIsInvalid + "\n",
		Names:   []string{"test_user", "test_user2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Admins: []string{},
				User:   []string{},
				Viewer: []string{},
			},
		},
	}, *res)
}

func testFailUserDoesNotExist(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "user": "test"})
	res := validateOrganisationCreate(ctx)
	assert.Equal(t, CreateOrganisationStruct{
		Message: fmt.Sprintf(generics.AccountDoesNotExistError, "test") + "\n",
		Names:   []string{"test_user", "test_user2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Admins: []string{},
				User:   []string{"test"},
				Viewer: []string{},
			},
		},
	}, *res)
}

func testFailAdminDoesNotExist(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "admins": "test"})
	res := validateOrganisationCreate(ctx)
	assert.Equal(t, CreateOrganisationStruct{
		Message: fmt.Sprintf(generics.AccountDoesNotExistError, "test") + "\n",
		Names:   []string{"test_user", "test_user2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Admins: []string{"test"},
				User:   []string{},
				Viewer: []string{},
			},
		},
	}, *res)
}

func testFailNoNames(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := validateOrganisationCreate(ctx)
	assert.Equal(t, CreateOrganisationStruct{
		Message: generics.NoMainGroupSubGroupOrNameProvided + "\n",
		Names:   []string{"test_user", "test_user2"},
		Organisation: database.Organisation{
			Info: database.OrganisationInfo{
				Admins: []string{},
				User:   []string{},
				Viewer: []string{},
			},
		},
	}, *res)
}

func setupCreateAccountsforOrgCreate(t *testing.T) {
	acc := database.Account{
		DisplayName: "test_user",
		Username:    "test_user",
		Password:    "test_user",
		Role:        database.User,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "test_user2", "test_user2"
	err = acc.CreateMe()
	assert.Nil(t, err)
}

func TestTestStructFillOrgCreateHandler(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupOrganisationCreate", setupOrganisationCreate)
	t.Run("testFillingOrgCreateStruct", testFillingOrgCreateStruct)
}

func testFillingOrgCreateStruct(t *testing.T) {
	result := getEmptyCreateOrgStruct()
	assert.Equal(t, CreateOrganisationStruct{ExistingMainGroup: []string{"a", "b"}, ExistingSubGroup: []string{"a"}, Names: []string{"test", "test2"}}, *result)
}

func setupOrganisationCreate(t *testing.T) {
	org := database.Organisation{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Status:    database.Public,
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.MainGroup = "b", "b"
	err = org.CreateMe()
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "test",
		Flair:       "",
		Username:    "test",
		Password:    "test",
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "test2", "test2"
	err = acc.CreateMe()
	assert.Nil(t, err)
}
