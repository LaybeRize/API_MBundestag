package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/htmlHandler"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrganisationLogic(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupAccountsForOrgLogic", testSetupAccountsForOrgLogic)
	t.Run("testCreateOrganisation", testCreateOrganisation)
	t.Run("testGetOrganisation", testGetOrganisation)
	t.Run("testUpdateOrganisation", testUpdateOrganisation)
	t.Run("testUpdateOnlyMemebers", testUpdateOnlyMemebers)
}

func testUpdateOnlyMemebers(t *testing.T) {
	org := Organsation{Name: "test_fail_orgLogic"}
	var msg htmlHandler.Message = ""
	var positiv = false
	org.ChangeOnlyMembers(&msg, &positiv)
	assert.Equal(t, false, positiv)
	assert.Equal(t, fmt.Sprintf(ErrorOrganisationNotFound, "test_fail_orgLogic")+"\n", string(msg))

	err := org.GetMe("test_orgLogic")
	assert.Nil(t, err)
	org.MainGroup = "tasdasd"
	org.Status = database.Hidden
	org.Member = []string{"test_OrgLogic", "test_OrgLogic2"}
	org.Admins = []string{"test"}

	msg = ""
	org.ChangeOnlyMembers(&msg, &positiv)
	assert.Equal(t, true, positiv)
	assert.Equal(t, SucessfulChangedOrganisation+"\n", msg)

	orgDB := database.Organisation{}
	err = orgDB.GetByName("test_orgLogic")
	assert.Nil(t, err)
	assert.Equal(t, "test_orgLogic2", orgDB.MainGroup)
	assert.Equal(t, "test2", orgDB.Flair.String)
	assert.Equal(t, database.Private, orgDB.Status)
	assert.Equal(t, 2, len(orgDB.Members))
	assert.Equal(t, 0, len(orgDB.Admins))
	assert.Equal(t, 1, len(orgDB.Accounts))
	assert.Equal(t, "test_OrgLogic", orgDB.Members[0].DisplayName)
	assert.Equal(t, "test_OrgLogic2", orgDB.Members[1].DisplayName)
	assert.Equal(t, "test_OrgLogic", orgDB.Accounts[0].DisplayName)
	acc := database.Account{}
	err = acc.GetByDisplayName("test_OrgLogic")
	assert.Nil(t, err)
	assert.Equal(t, "test2", acc.Flair)
	err = acc.GetByDisplayNameWithParent("test_OrgLogic2")
	assert.Nil(t, err)
	assert.Equal(t, "test2", acc.Flair)
	err = acc.GetByDisplayNameWithParent("test_OrgLogic3")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
}

func testUpdateOrganisation(t *testing.T) {
	org := Organsation{Name: "test_fail_orgLogic"}
	var msg htmlHandler.Message = ""
	var positiv = false
	org.ChangeMe(&msg, &positiv)
	assert.Equal(t, false, positiv)
	assert.Equal(t, fmt.Sprintf(ErrorOrganisationNotFound, "test_fail_orgLogic")+"\n", string(msg))

	err := org.GetMe("test_orgLogic")
	assert.Nil(t, err)

	org.MainGroup = "test_orgLogic2"
	org.Flair = "test2"
	org.Status = database.Private
	org.Member = []string{"test_OrgLogic2", "test_OrgLogic3"}
	org.Admins = []string{}

	msg = ""
	org.ChangeMe(&msg, &positiv)
	assert.Equal(t, true, positiv)
	assert.Equal(t, SucessfulChangedOrganisation+"\n", msg)

	orgDB := database.Organisation{}
	err = orgDB.GetByName("test_orgLogic")
	assert.Nil(t, err)
	assert.Equal(t, "test_orgLogic2", orgDB.MainGroup)
	assert.Equal(t, "test2", orgDB.Flair.String)
	assert.Equal(t, database.Private, orgDB.Status)
	assert.Equal(t, 2, len(orgDB.Members))
	assert.Equal(t, 0, len(orgDB.Admins))
	assert.Equal(t, 2, len(orgDB.Accounts))
	assert.Equal(t, "test_OrgLogic2", orgDB.Members[0].DisplayName)
	assert.Equal(t, "test_OrgLogic3", orgDB.Members[1].DisplayName)
	assert.Equal(t, "test_OrgLogic", orgDB.Accounts[0].DisplayName)
	assert.Equal(t, "test_OrgLogic3", orgDB.Accounts[1].DisplayName)
	acc := database.Account{}
	err = acc.GetByDisplayName("test_OrgLogic")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
	err = acc.GetByDisplayNameWithParent("test_OrgLogic2")
	assert.Nil(t, err)
	assert.Equal(t, "test2", acc.Flair)
	err = acc.GetByDisplayNameWithParent("test_OrgLogic3")
	assert.Nil(t, err)
	assert.Equal(t, "test2", acc.Flair)
}

func testGetOrganisation(t *testing.T) {
	expected := Organsation{
		Name:      "test_orgLogic",
		MainGroup: "test_orgLogic",
		SubGroup:  "test_orgLogic",
		Flair:     "test",
		Status:    database.Public,
		Member:    []string{"test_OrgLogic"},
		Admins:    []string{"test_OrgLogic2"},
	}
	org := Organsation{}
	err := org.GetMe("test_orgLogic")
	assert.Nil(t, err)
	assert.Equal(t, expected, org)
}

func testCreateOrganisation(t *testing.T) {
	org := Organsation{
		Name:      "test_orgLogic",
		MainGroup: "test_orgLogic",
		SubGroup:  "test_orgLogic",
		Flair:     "test",
		Status:    database.Public,
		Member:    []string{"fail_me_now"},
		Admins:    []string{"fail_me_now"},
	}
	var msg htmlHandler.Message = ""
	var positiv = false

	org.CreateMe(&msg, &positiv)
	assert.Equal(t, false, positiv)
	assert.Equal(t, fmt.Sprintf(AccountDoesNotExistError, "fail_me_now")+"\n", string(msg))
	org.Member = []string{"test_OrgLogic"}
	msg = ""
	org.CreateMe(&msg, &positiv)
	assert.Equal(t, false, positiv)
	assert.Equal(t, fmt.Sprintf(AccountDoesNotExistError, "fail_me_now")+"\n", string(msg))
	org.Admins = []string{"test_OrgLogic2"}
	msg = ""
	org.CreateMe(&msg, &positiv)
	assert.Equal(t, true, positiv)
	assert.Equal(t, OrganisationSuccessfulCreated+"\n", msg)

	orgDB := database.Organisation{}
	err := orgDB.GetByName("test_orgLogic")
	assert.Nil(t, err)
	assert.Equal(t, org.Name, orgDB.Name)
	assert.Equal(t, org.MainGroup, orgDB.MainGroup)
	assert.Equal(t, org.SubGroup, orgDB.SubGroup)
	assert.Equal(t, org.Flair, orgDB.Flair.String)
	assert.Equal(t, org.Status, orgDB.Status)
	assert.Equal(t, 1, len(orgDB.Members))
	assert.Equal(t, 1, len(orgDB.Admins))
	assert.Equal(t, org.Member[0], orgDB.Members[0].DisplayName)
	assert.Equal(t, org.Admins[0], orgDB.Admins[0].DisplayName)
	assert.Equal(t, 1, len(orgDB.Accounts))
	assert.Equal(t, org.Member[0], orgDB.Accounts[0].DisplayName)
	acc := database.Account{}
	err = acc.GetByDisplayName("test_OrgLogic")
	assert.Nil(t, err)
	assert.Equal(t, "test", acc.Flair)
	err = acc.GetByDisplayNameWithParent("test_OrgLogic2")
	assert.Nil(t, err)
	assert.Equal(t, "test", acc.Flair)
	err = acc.GetByDisplayNameWithParent("test_OrgLogic3")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
}

func testSetupAccountsForOrgLogic(t *testing.T) {
	acc := database.Account{
		DisplayName: "test_OrgLogic",
		Username:    "test_OrgLogic",
		Password:    "test",
		Role:        database.User,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	id := acc.ID
	acc = database.Account{
		DisplayName: "test_OrgLogic2",
		Username:    "test_OrgLogic2",
		Password:    "XXXX",
		Role:        database.User,
		Linked:      sql.NullInt64{Valid: true, Int64: id},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "test_OrgLogic3",
		Username:    "test_OrgLogic3",
		Password:    "XXXX",
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
}
