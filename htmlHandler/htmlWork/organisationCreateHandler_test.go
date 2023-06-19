package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestOrganisationCreateHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupOrganisationCreatePage", testSetupOrganisationCreatePage)
	t.Run("testGetOrganisationCreatePage", testGetOrganisationCreatePage)
	t.Run("testPostOrganisationCreatePage", testPostOrganisationCreatePage)
}

func testPostOrganisationCreatePage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testCreateHandlerError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostCreateOrganisationPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testCreateHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	PostCreateOrganisationPage(ctx)
	assert.Equal(t, "TestCreateOrganisation|{     [] []}|"+string(generics.NoMainGroupSubGroupOrNameProvided)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testOrgCreateHTML", "mainGroup": "testOrgCreateHTML", "subGroup": "testOrgCreateHTML"})
	PostCreateOrganisationPage(ctx)
	assert.Equal(t, "TestCreateOrganisation|{testOrgCreateHTML testOrgCreateHTML testOrgCreateHTML   [] []}|"+string(generics.StatusIsInvalid)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testOrgCreateHTML", "mainGroup": "testOrgCreateHTML", "subGroup": "testOrgCreateHTML", "status": string(database.Public)})
	PostCreateOrganisationPage(ctx)
	assert.Equal(t, "TestCreateOrganisation|{testOrgCreateHTML testOrgCreateHTML testOrgCreateHTML  public [] []}|"+string(dataLogic.OrganisationSuccessfulCreated)+"\n|true", w.Body.String())

	org := database.Organisation{}
	err = org.GetByName("testOrgCreateHTML")
	assert.Nil(t, err)
	assert.Equal(t, "testOrgCreateHTML", org.Name)
	assert.Equal(t, "testOrgCreateHTML", org.MainGroup)
	assert.Equal(t, "testOrgCreateHTML", org.SubGroup)
	assert.Equal(t, database.Public, org.Status)
	assert.Equal(t, 0, len(org.Members))
	assert.Equal(t, 0, len(org.Admins))
}

func testGetOrganisationCreatePage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testCreateHandlerError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetCreateOrganisationPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testCreateHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUser(t, acc)
	GetCreateOrganisationPage(ctx)
	assert.Equal(t, "TestCreateOrganisation|{    "+string(database.Public)+" [] []}||false", w.Body.String())
}

func testSetupOrganisationCreatePage(t *testing.T) {
	Setup()
	htmlBasics.Setup()

	acc := database.Account{
		DisplayName: "testCreateHandlerError",
		Username:    "testCreateHandlerError",
		Role:        database.MediaAdmin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testCreateHandler",
		Username:    "testCreateHandler",
		Role:        database.Admin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestCreateOrganisation|{{.Page.Organisation}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "createOrganisation")
}
