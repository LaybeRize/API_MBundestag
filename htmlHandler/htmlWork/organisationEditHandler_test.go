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
	"testing"
)

func TestOrgansationEditHandlerPage(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupOrganisationEdit", setupOrganisationEdit)
	t.Run("additionSetupOrgEdit", additionSetupOrgEdit)
	t.Run("setupPageEditOrganisation", setupPageEditOrganisation)
	t.Run("testGetEditOrganisationPage", testGetEditOrganisationPage)
	t.Run("testPostSearchEditOrganisationPage", testPostSearchEditOrganisationPage)
	t.Run("testPostEditOrganisationPage", testPostEditOrganisationPage)
}

func testPostEditOrganisationPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("test_admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "a", "mainGroup": "b", "admins": "press2", "user": "press", "subGroup": "a", "status": string(database.Secret)})
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "editOrganisation [a b] [a b] [a] a b a  {[press2] [press] [test2 test]} "+string(database.Secret)+generics.SuccessFullChangeOrg+"\n", w.Body.String())
}

func testPostSearchEditOrganisationPage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("test_admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "a"})
	ctx.Request.URL.RawQuery = "type=search"
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "editOrganisation [a b] [a b] [a] a a a  {[] [] []} "+string(database.Public)+generics.SuccessFullFindOrg+"\n", w.Body.String())
}

func testGetEditOrganisationPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByDisplayName("test_admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "editOrganisation [a b] [a b] [a]     {[] [] []} "+string(database.Public), w.Body.String())
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "org=a"
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "editOrganisation [a b] [a b] [a] a a a  {[] [] []} "+string(database.Public), w.Body.String())
}

func setupPageEditOrganisation(t *testing.T) {
	acc := database.Account{
		DisplayName: "test_admin",
		Username:    "test_admin",
		Role:        database.Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.OrgNames}} {{.Page.ExistingMainGroup}} {{.Page.ExistingSubGroup}} {{.Page.Organisation.Name}} {{.Page.Organisation.MainGroup}} {{.Page.Organisation.SubGroup}} {{.Page.Organisation.Flair.String}} {{.Page.Organisation.Info}} {{.Page.Organisation.Status}}{{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":            temp,
			"editOrganisation": temp2,
		},
	}
}

func TestValidateEditOrganisation(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupOrganisationEdit", setupOrganisationEdit)
	t.Run("additionSetupOrgEdit", additionSetupOrgEdit)
	t.Run("testNoMainGroupSubGroupOrNameProvided", testNoMainGroupSubGroupOrNameProvided)
	t.Run("testOrgEditNonExistantElement", testOrgEditNonExistantElement)
	t.Run("testAdminAccountDoesNotExistError", testAdminAccountDoesNotExistError)
	t.Run("testUserAccountDoesNotExistError", testUserAccountDoesNotExistError)
	t.Run("testStatusIsInvalid", testStatusIsInvalid)
	t.Run("testSuccessFullChangeOrg", testSuccessFullChangeOrg)
	t.Run("testSuccessfulFullChangeOrg", testSuccessfulFullChangeOrg)
	t.Run("testHiddingChangeOrg", testHiddingChangeOrg)
}

func testHiddingChangeOrg(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"flair": "b", "name": "a", "mainGroup": "b", "admins": "press2", "user": "press", "subGroup": "a", "status": string(database.Hidden)})
	ctx.Request.PostForm.Add("admins", "press")
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.SuccessFullChangeOrg + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "b",
			SubGroup:  "a",
			Status:    database.Hidden,
			Info: database.OrganisationInfo{
				Viewer: []string{},
				User:   []string{},
				Admins: []string{},
			}},
	}, *res)
	acc := database.Account{}
	err := acc.GetByUserName("press")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
	err = acc.GetByUserName("press2")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
}

func testSuccessfulFullChangeOrg(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"flair": "b", "name": "a", "mainGroup": "b", "admins": "press2", "user": "press", "subGroup": "a", "status": string(database.Public)})
	ctx.Request.PostForm.Add("admins", "press")
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.SuccessFullChangeOrg + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "b",
			SubGroup:  "a",
			Flair:     sql.NullString{Valid: true, String: "b"},
			Status:    database.Public,
			Info: database.OrganisationInfo{
				Viewer: []string{"test2", "test"},
				User:   []string{"press"},
				Admins: []string{"press2"},
			}},
	}, *res)
	acc := database.Account{}
	err := acc.GetByUserName("press")
	assert.Nil(t, err)
	assert.Equal(t, "b", acc.Flair)
	err = acc.GetByUserName("press2")
	assert.Nil(t, err)
	assert.Equal(t, "b", acc.Flair)
}

func testSuccessFullChangeOrg(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "status": string(database.Secret)})
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.SuccessFullChangeOrg + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Status:    database.Secret,
			Info: database.OrganisationInfo{
				Viewer: []string{},
				User:   []string{},
				Admins: []string{},
			}},
	}, *res)
}

func testStatusIsInvalid(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a"})
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.StatusIsInvalid + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Viewer: []string{},
				User:   []string{},
				Admins: []string{},
			}},
	}, *res)
}

func testUserAccountDoesNotExistError(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "user": "lol"})
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           fmt.Sprintf(generics.AccountDoesNotExistError, "lol") + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Viewer: []string{},
				User:   []string{"lol"},
				Admins: []string{},
			}},
	}, *res)
}

func testAdminAccountDoesNotExistError(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a", "mainGroup": "a", "subGroup": "a", "admins": "lol"})
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           fmt.Sprintf(generics.AccountDoesNotExistError, "lol") + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "a",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Viewer: []string{},
				User:   []string{},
				Admins: []string{"lol"},
			}},
	}, *res)
}

func testOrgEditNonExistantElement(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "c", "mainGroup": "a", "subGroup": "a"})
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.OrgEditNonExistantElement + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{
			Name:      "c",
			MainGroup: "a",
			SubGroup:  "a",
			Info: database.OrganisationInfo{
				Viewer: []string{},
				User:   []string{},
				Admins: []string{},
			}},
	}, *res)
}

func testNoMainGroupSubGroupOrNameProvided(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := validateOrganisationEdit(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.NoMainGroupSubGroupOrNameProvided + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"press", "press2", "test", "test2"},
		Organisation: database.Organisation{Info: database.OrganisationInfo{
			Viewer: []string{},
			User:   []string{},
			Admins: []string{},
		}},
	}, *res)
}

func additionSetupOrgEdit(t *testing.T) {
	acc := database.Account{}
	acc.DisplayName, acc.Username = "press", "press"
	acc.Linked.Valid = true
	acc.Role, acc.Linked.Int64 = database.PressAccount, 1
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username, acc.Linked.Int64 = "press2", "press2", 2
	err = acc.CreateMe()
	assert.Nil(t, err)
}

func TestValidateGetOrgStruct(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupOrganisationEdit", setupOrganisationEdit)
	t.Run("testFailGetOrganisationStruct", testFailGetOrganisationStruct)
	t.Run("testFailGetOrganisationStruct", testSuccessGetOrganisationStruct)
}

func testSuccessGetOrganisationStruct(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"name": "a"})
	org := database.Organisation{}
	err := org.GetByName("a")
	assert.Nil(t, err)
	res := vaildateOrganisationSearch(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.SuccessFullFindOrg + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"test", "test2"},
		Organisation:      org,
	}, *res)
}

func testFailGetOrganisationStruct(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	res := vaildateOrganisationSearch(ctx)
	assert.Equal(t, EditOrganisationStruct{
		Message:           generics.OrgFindingError + "\n",
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"test", "test2"},
		Organisation: database.Organisation{
			Status: database.Public,
			Info: database.OrganisationInfo{
				Admins: []string{},
				User:   []string{},
			},
		},
	}, *res)
}

func TestStructEditOrganisation(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestOrganisationsDB()

	t.Run("setupOrganisationEdit", setupOrganisationEdit)
	t.Run("testFillingOrgEditStruct", testFillingOrgEditStruct)
}

func testFillingOrgEditStruct(t *testing.T) {
	res := getEmptyEditOrgStruct()
	assert.Equal(t, EditOrganisationStruct{
		Organisation:      database.Organisation{},
		OrgNames:          []string{"a", "b"},
		ExistingMainGroup: []string{"a", "b"},
		ExistingSubGroup:  []string{"a"},
		Names:             []string{"test", "test2"},
		Message:           "",
	}, *res)
}

func setupOrganisationEdit(t *testing.T) {
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
	acc.DisplayName, acc.Username, acc.Suspended = "sus", "sus", true
	err = acc.CreateMe()
	assert.Nil(t, err)
	err = acc.SaveChanges()
	assert.Nil(t, err)
}
