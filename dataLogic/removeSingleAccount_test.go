package dataLogic

import (
	"API_MBundestag/database"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var singleAccountTests database.Account
var singleAccountTestPress database.Account

func TestSingleRemoval(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupOrgsAndTitlesAndAccounts", testSetupOrgsAndTitlesAndAccounts)
	t.Run("testRemoveFromTitle", testRemoveFromTitle)
	t.Run("testRemoveFromOrg", testRemoveFromOrg)
}

func testRemoveFromOrg(t *testing.T) {
	err := RemoveSingleAccountFromOrganisations(&singleAccountTests)
	assert.Nil(t, err)
	org := database.Organisation{}
	err = org.GetByName("test_singleAccount")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, 1, len(org.Admins))
	assert.Equal(t, 2, len(org.Accounts))
	assert.Equal(t, "head_admin", org.Members[0].DisplayName)
	assert.Equal(t, "singleAccount_test_press", org.Admins[0].DisplayName)
	assert.Equal(t, "head_admin", org.Accounts[0].DisplayName)
	assert.Equal(t, "singleAccount_test", org.Accounts[1].DisplayName)
	err = org.GetByName("test_singleAccount2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, 1, len(org.Admins))
	assert.Equal(t, 2, len(org.Accounts))
	assert.Equal(t, "head_admin", org.Members[0].DisplayName)
	assert.Equal(t, "singleAccount_test_press", org.Admins[0].DisplayName)
	assert.Equal(t, "head_admin", org.Accounts[0].DisplayName)
	assert.Equal(t, "singleAccount_test", org.Accounts[1].DisplayName)
	err = org.GetByName("test_singleAccount3")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, 0, len(org.Admins))
	assert.Equal(t, 1, len(org.Accounts))
	assert.Equal(t, "head_admin", org.Members[0].DisplayName)
	assert.Equal(t, "head_admin", org.Accounts[0].DisplayName)

	err = singleAccountTests.GetByUserName(singleAccountTests.Username)
	assert.Nil(t, err)
	assert.Equal(t, "", singleAccountTests.Flair)

	err = RemoveSingleAccountFromOrganisations(&singleAccountTestPress)
	assert.Nil(t, err)
	err = org.GetByName("test_singleAccount")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, 0, len(org.Admins))
	assert.Equal(t, 1, len(org.Accounts))
	assert.Equal(t, "head_admin", org.Members[0].DisplayName)
	assert.Equal(t, "head_admin", org.Accounts[0].DisplayName)
	err = org.GetByName("test_singleAccount2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, 0, len(org.Admins))
	assert.Equal(t, 1, len(org.Accounts))
	assert.Equal(t, "head_admin", org.Members[0].DisplayName)
	assert.Equal(t, "head_admin", org.Accounts[0].DisplayName)
	err = org.GetByName("test_singleAccount3")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(org.Members))
	assert.Equal(t, 0, len(org.Admins))
	assert.Equal(t, 1, len(org.Accounts))
	assert.Equal(t, "head_admin", org.Members[0].DisplayName)
	assert.Equal(t, "head_admin", org.Accounts[0].DisplayName)

	err = singleAccountTestPress.GetByUserName(singleAccountTestPress.Username)
	assert.Nil(t, err)
	assert.Equal(t, "", singleAccountTestPress.Flair)
}

func testRemoveFromTitle(t *testing.T) {
	err := RemoveSingleAccountFromTitles(&singleAccountTests)
	assert.Nil(t, err)

	title := database.Title{}
	err = title.GetByName("remove_title_test")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(title.Holder))
	assert.Equal(t, "head_admin", title.Holder[0].DisplayName)
	err = title.GetByName("remove_title_test2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(title.Holder))
	assert.Equal(t, "singleAccount_test_press", title.Holder[0].DisplayName)

	err = singleAccountTests.GetByUserName(singleAccountTests.Username)
	assert.Nil(t, err)
	assert.Equal(t, "test_acc_2", singleAccountTests.Flair)

	err = RemoveSingleAccountFromTitles(&singleAccountTestPress)
	assert.Nil(t, err)

	err = title.GetByName("remove_title_test")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(title.Holder))
	assert.Equal(t, "head_admin", title.Holder[0].DisplayName)
	err = title.GetByName("remove_title_test2")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(title.Holder))

	err = singleAccountTestPress.GetByUserName(singleAccountTestPress.Username)
	assert.Nil(t, err)
	assert.Equal(t, "test_acc_1, test_acc_2", singleAccountTestPress.Flair)
}

func testSetupOrgsAndTitlesAndAccounts(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("head_admin")
	assert.Nil(t, err)

	singleAccountTests = database.Account{
		DisplayName: "singleAccount_test",
		Username:    "singleAccount_test",
		Password:    "test",
		Flair:       "test_acc_2, test_title2, test_title3",
		Role:        database.User,
	}
	err = singleAccountTests.CreateMe()
	assert.Nil(t, err)

	singleAccountTestPress = database.Account{
		DisplayName: "singleAccount_test_press",
		Username:    "singleAccount_test_press",
		Password:    "XXXX",
		Flair:       "test_acc_1, test_acc_2, test_title3",
		Role:        database.PressAccount,
		Linked:      sql.NullInt64{Valid: true, Int64: singleAccountTests.ID},
	}
	err = singleAccountTestPress.CreateMe()
	assert.Nil(t, err)

	title := database.Title{
		Name:      "remove_title_test",
		MainGroup: "remove_title_test",
		SubGroup:  "remove_title_test",
		Flair:     sql.NullString{Valid: true, String: "test_title2"},
		Holder:    []database.Account{acc, singleAccountTests},
	}
	err = title.CreateMe()
	assert.Nil(t, err)

	title = database.Title{
		Name:      "remove_title_test2",
		MainGroup: "remove_title_test2",
		SubGroup:  "remove_title_test2",
		Flair:     sql.NullString{Valid: true, String: "test_title3"},
		Holder:    []database.Account{singleAccountTests, singleAccountTestPress},
	}
	err = title.CreateMe()
	assert.Nil(t, err)

	org := database.Organisation{
		Name:      "test_singleAccount",
		MainGroup: "test_singleAccount",
		SubGroup:  "test_singleAccount",
		Flair: sql.NullString{
			String: "test_acc_1",
			Valid:  true,
		},
		Status:   database.Public,
		Members:  []database.Account{acc},
		Admins:   []database.Account{singleAccountTestPress},
		Accounts: []database.Account{acc, singleAccountTests},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = database.Organisation{
		Name:      "test_singleAccount2",
		MainGroup: "test_singleAccount2",
		SubGroup:  "test_singleAccount2",
		Flair: sql.NullString{
			String: "test_acc_2",
			Valid:  true,
		},
		Status:   database.Private,
		Members:  []database.Account{acc, singleAccountTests},
		Admins:   []database.Account{singleAccountTestPress},
		Accounts: []database.Account{acc, singleAccountTests},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = database.Organisation{
		Name:      "test_singleAccount3",
		MainGroup: "test_singleAccount3",
		SubGroup:  "test_singleAccount3",
		Flair:     sql.NullString{},
		Status:    database.Secret,
		Members:   []database.Account{acc},
		Admins:    []database.Account{singleAccountTests},
		Accounts:  []database.Account{acc, singleAccountTests},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
}
