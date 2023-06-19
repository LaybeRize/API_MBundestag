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

func TestOrganisationUserHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupOrganisationUserPage", testSetupOrganisationUserPage)
	t.Run("testGetOrganisationUserPage", testGetOrganisationUserPage)
	t.Run("testPostOrganisationUserPage", testPostOrganisationUserPage)
}

func testPostOrganisationUserPage(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	PostOrganisationUserHandler(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("testEditUserOrgHandlerError")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "name=testEditUserOrgHandler")
	GetOrganisationUserHandler(ctx)
	assert.Equal(t, "TestError|"+generics.CouldNotFindOrganisation, w.Body.String())

	err = acc.GetByUserName("testEditOrgUserHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "bakjsdöjasd"})
	PostOrganisationUserHandler(ctx)
	assert.Equal(t, "TestError|"+generics.CouldNotFindOrganisation, w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testEditUserOrgHandler", "user": "testEditOrgUserHandler"})
	PostOrganisationUserHandler(ctx)
	assert.Equal(t, "TestUserEditOrganisation|testEditUserOrgHandler|[]|"+string(dataLogic.SucessfulChangedOrganisation)+"\n|true", w.Body.String())

	org := database.Organisation{}
	err = org.GetByName("testEditUserOrgHandler")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(org.Members))

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testEditUserOrgHandler", "user": "testEditUserOrgHandlerError"})
	PostOrganisationUserHandler(ctx)
	assert.Equal(t, "TestUserEditOrganisation|testEditUserOrgHandler|[testEditUserOrgHandlerError]|"+string(dataLogic.SucessfulChangedOrganisation)+"\n|true", w.Body.String())

	err = org.GetByName("testEditUserOrgHandler")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, "testEditUserOrgHandlerError", org.Members[0].DisplayName)
}

func testGetOrganisationUserPage(t *testing.T) {
	w, ctx := hHa.GetEmptyContext(t)
	GetOrganisationUserHandler(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("testEditUserOrgHandlerError")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "name=testEditUserOrgHandler")
	GetOrganisationUserHandler(ctx)
	assert.Equal(t, "TestError|"+generics.CouldNotFindOrganisation, w.Body.String())

	err = acc.GetByUserName("testEditOrgUserHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "name=bakjsdöjasd")
	GetOrganisationUserHandler(ctx)
	assert.Equal(t, "TestError|"+generics.CouldNotFindOrganisation, w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "name=testEditUserOrgHandler")
	GetOrganisationUserHandler(ctx)
	assert.Equal(t, "TestUserEditOrganisation|testEditUserOrgHandler|[]|"+string(generics.SuccessFullFindOrg)+"|true", w.Body.String())
}

func testSetupOrganisationUserPage(t *testing.T) {
	Setup()
	htmlBasics.Setup()

	acc := database.Account{
		DisplayName: "testEditUserOrgHandlerError",
		Username:    "testEditUserOrgHandlerError",
		Role:        database.MediaAdmin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testEditOrgUserHandler",
		Username:    "testEditOrgUserHandler",
		Flair:       "testEditHandler",
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	org := database.Organisation{
		Name:      "testEditUserOrgHandler",
		MainGroup: "testEditHandler",
		SubGroup:  "testEditHandler",
		Flair:     sql.NullString{Valid: true, String: "testEditHandler"},
		Status:    database.Public,
		Admins:    []database.Account{acc},
	}
	err = org.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestUserEditOrganisation|{{.Page.OrganisationName}}|{{.Page.User}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "editUserOrganisation")
}
