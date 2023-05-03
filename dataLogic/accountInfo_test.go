package dataLogic

import (
	"API_MBundestag/database"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testaccountinfoAcc database.Account
var testaccountinfoAcc2 database.Account

func TestAccountInfo(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupAccountInfo", testSetupAccountInfo)
	t.Run("testGetTitleList", testGetTitleList)
	t.Run("testGetOrganisationList", testGetOrganisationList)
}

func testGetTitleList(t *testing.T) {
	err, str := GetTitelList(testaccountinfoAcc.ID)
	assert.Nil(t, err)
	assert.Equal(t, "", str)
	err, str = GetTitelList(testaccountinfoAcc2.ID)
	assert.Nil(t, err)
	assert.Equal(t, "", str)

	title := database.Title{
		Name:   "accInfo_title",
		Holder: []database.Account{testaccountinfoAcc, testaccountinfoAcc2},
	}
	err = title.CreateMe()
	assert.Nil(t, err)
	title.Name = "accInfo_title2"
	err = title.CreateMe()
	assert.Nil(t, err)
	title.Name, title.Holder = "accInfo_title3", []database.Account{testaccountinfoAcc}
	err = title.CreateMe()
	assert.Nil(t, err)

	err, str = GetTitelList(testaccountinfoAcc.ID)
	assert.Nil(t, err)
	assert.Equal(t, "accInfo_title, accInfo_title2, accInfo_title3", str)
	err, str = GetTitelList(testaccountinfoAcc2.ID)
	assert.Nil(t, err)
	assert.Equal(t, "accInfo_title, accInfo_title2", str)
}

func testGetOrganisationList(t *testing.T) {
	err, str := GetOrganisationList(testaccountinfoAcc.ID)
	assert.Nil(t, err)
	assert.Equal(t, "", str)
	err, str = GetOrganisationList(testaccountinfoAcc2.ID)
	assert.Nil(t, err)
	assert.Equal(t, "", str)

	org := database.Organisation{
		Name:      "accInfo_org",
		MainGroup: "accInfo_org",
		SubGroup:  "accInfo_org",
		Flair:     sql.NullString{},
		Status:    database.Public,
		Members:   []database.Account{testaccountinfoAcc},
		Admins:    []database.Account{testaccountinfoAcc2},
		Accounts:  []database.Account{},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = database.Organisation{
		Name:      "accInfo_org2",
		MainGroup: "accInfo_org2",
		SubGroup:  "accInfo_org2",
		Flair:     sql.NullString{},
		Status:    database.Public,
		Members:   []database.Account{testaccountinfoAcc2},
		Admins:    []database.Account{},
		Accounts:  []database.Account{testaccountinfoAcc},
	}
	err = org.CreateMe()
	assert.Nil(t, err)

	err, str = GetOrganisationList(testaccountinfoAcc.ID)
	assert.Nil(t, err)
	assert.Equal(t, "accInfo_org", str)
	err, str = GetOrganisationList(testaccountinfoAcc2.ID)
	assert.Nil(t, err)
	assert.Equal(t, "accInfo_org, accInfo_org2", str)
}

func testSetupAccountInfo(t *testing.T) {
	testaccountinfoAcc = database.Account{
		DisplayName: "accInfo1",
		Username:    "accInfo1",
		Password:    "test",
		Role:        database.User,
	}
	testaccountinfoAcc2 = testaccountinfoAcc
	testaccountinfoAcc2.DisplayName = "accInfo2"
	testaccountinfoAcc2.Username = "accInfo2"
	err := testaccountinfoAcc.CreateMe()
	assert.Nil(t, err)
	err = testaccountinfoAcc2.CreateMe()
	assert.Nil(t, err)
}
