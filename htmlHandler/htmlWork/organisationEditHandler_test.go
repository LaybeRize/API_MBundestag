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

func TestOrganisationEditHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupOrganisationEditHandler", testSetupOrganisationEditHandler)
	t.Run("testGetOrganisationEditHandler", testGetOrganisationEditHandler)
	t.Run("testPostSearchOrganisationEditHandler", testPostSearchOrganisationEditHandler)
	t.Run("testPostOrganisationEditHandler", testPostOrganisationEditHandler)
}

func testPostOrganisationEditHandler(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditHandler")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{     [] []}|"+string(generics.NoMainGroupSubGroupOrNameProvided)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "asvasdasdasd", "mainGroup": "testOrgCreateHTML", "subGroup": "testOrgCreateHTML"})
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{asvasdasdasd testOrgCreateHTML testOrgCreateHTML   [] []}|"+string(generics.OrgEditNonExistantElement)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testEditHandler", "mainGroup": "testOrgCreateHTML", "subGroup": "testOrgCreateHTML", "flair": "flair"})
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{testEditHandler testOrgCreateHTML testOrgCreateHTML flair  [] []}|"+string(generics.StatusIsInvalid)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithForm(t, acc, map[string]string{"name": "testEditHandler", "mainGroup": "testOrgCreateHTML", "subGroup": "testOrgCreateHTML", "status": string(database.Private), "flair": "flair"})
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{testEditHandler testOrgCreateHTML testOrgCreateHTML flair "+string(database.Private)+" [] []}|"+string(dataLogic.SucessfulChangedOrganisation)+"\n|true", w.Body.String())

	org := database.Organisation{}
	err = org.GetByName("testEditHandler")
	assert.Nil(t, err)
	assert.Equal(t, "testEditHandler", org.Name)
	assert.Equal(t, "testOrgCreateHTML", org.MainGroup)
	assert.Equal(t, "testOrgCreateHTML", org.SubGroup)
	assert.Equal(t, "flair", org.Flair.String)
	assert.Equal(t, database.Private, org.Status)
	assert.Equal(t, 0, len(org.Members))
	assert.Equal(t, 0, len(org.Admins))
}

func testPostSearchOrganisationEditHandler(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditHandlerError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testEditHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "type=search")
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{    "+string(database.Public)+" [] []}|"+string(generics.OrgFindingError)+"\n|false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{"name": "testEditHandler"}, "type=search")
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{testEditHandler testEditHandler testEditHandler testEditHandler "+string(database.Public)+" [] []}|"+string(generics.SuccessFullFindOrg)+"\n|true", w.Body.String())
}

func testGetOrganisationEditHandler(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditHandlerError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

	err = acc.GetByUserName("testEditHandler")
	assert.Nil(t, err)
	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "org=basdaewaw1234eadsd")
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{    "+string(database.Public)+" [] []}||false", w.Body.String())

	w, ctx = hHa.GetContextSetForUserWithFormAndQuery(t, acc, map[string]string{}, "org=testEditHandler")
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "TestEditOrganisation|{testEditHandler testEditHandler testEditHandler testEditHandler "+string(database.Public)+" [] []}||false", w.Body.String())
}

func testSetupOrganisationEditHandler(t *testing.T) {
	Setup()
	htmlBasics.Setup()

	org := database.Organisation{
		Name:      "testEditHandler",
		MainGroup: "testEditHandler",
		SubGroup:  "testEditHandler",
		Flair:     sql.NullString{Valid: true, String: "testEditHandler"},
		Status:    database.Public,
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "testEditHandlerError",
		Username:    "testEditHandlerError",
		Role:        database.MediaAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testEditHandler",
		Username:    "testEditHandler",
		Role:        database.Admin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestEditOrganisation|{{.Page.Organisation}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "editOrganisation")
}
