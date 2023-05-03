package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var orgAcc1ID = int64(0)
var orgAcc2ID = int64(0)

func TestOrganisationManagement(t *testing.T) {
	TestSetup()
	t.Run("testSetupAccountsOrganisation", testSetupAccountsOrganisation)
	t.Run("testOrgansationCreate", testOrgansationCreate)
	t.Run("testOrganisationChange", testOrganisationChange)
	t.Run("testGetByAccounts", testGetByAccounts)
	t.Run("testSetupOrganisations", testSetupOrganisations)
	t.Run("testOrgUserLists", testOrgUserLists)
	t.Run("testOrgLists", testOrgLists)
	t.Run("testOrgGroups", testOrgGroups)
}

func testOrgGroups(t *testing.T) {
	testExists := false
	test2Exists := false
	list := OrganisationList{}
	err := list.GetAllSubGroups()
	assert.Nil(t, err)
	for _, org := range list {
		if org.SubGroup == "test" && testExists {
			t.Fail()
		} else if org.SubGroup == "test" {
			testExists = true
		}

		if org.SubGroup == "test2" && test2Exists {
			t.Fail()
		} else if org.SubGroup == "test2" {
			test2Exists = true
		}
	}
	assert.True(t, testExists)
	assert.True(t, test2Exists)
	testExists = false
	test2Exists = false
	err = list.GetAllMainGroups()
	assert.Nil(t, err)
	for _, org := range list {
		if org.MainGroup == "test" && testExists {
			t.Fail()
		} else if org.MainGroup == "test" {
			testExists = true
		}

		if org.MainGroup == "test2" && test2Exists {
			t.Fail()
		} else if org.MainGroup == "test2" {
			test2Exists = true
		}
	}
	assert.True(t, testExists)
	assert.True(t, test2Exists)
}

func testOrgLists(t *testing.T) {
	list := OrganisationList{}
	err := list.GetAllVisable()
	assert.Nil(t, err)
	counter := 0
	for _, org := range list {
		switch org.Name {
		case "o_test_org", "o_test_org2", "o_test_org3":
			counter++
		}
	}
	assert.Equal(t, 3, counter)
	err = list.GetAllInvisable()
	assert.Nil(t, err)
	exists := false
	for _, org := range list {
		if org.Name == "o_test_org_hidden" {
			exists = true
		}
	}
	assert.True(t, exists)
}

func testOrgUserLists(t *testing.T) {
	list := OrganisationList{}
	err := list.GetAllVisibleFor(orgAcc1ID)
	assert.Nil(t, err)
	counter := 0
	for _, org := range list {
		switch org.Name {
		case "o_test_org", "o_test_org2", "o_test_org3":
			counter++
		}
	}
	assert.Equal(t, 3, counter)
	err = list.GetAllVisibleFor(orgAcc2ID)
	assert.Nil(t, err)
	counter = 0
	for _, org := range list {
		switch org.Name {
		case "o_test_org", "o_test_org2", "o_test_org3":
			counter++
		}
	}
	assert.Equal(t, 2, counter)
	err = list.GetAllPartOf(orgAcc1ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "o_test_org", list[0].Name)
	assert.Equal(t, "o_test_org2", list[1].Name)
	err = list.GetAllPartOf(orgAcc2ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "o_test_org", list[0].Name)
	assert.Equal(t, "o_test_org3", list[1].Name)
}

func testSetupOrganisations(t *testing.T) {
	acc1 := Account{}
	err := acc1.GetByUserName("test_org1")
	assert.Nil(t, err)
	orgAcc1ID = acc1.ID
	acc2 := Account{}
	err = acc2.GetByUserName("test_org2")
	assert.Nil(t, err)
	orgAcc2ID = acc2.ID
	acc3 := Account{}
	err = acc3.GetByUserName("test_org3")
	assert.Nil(t, err)
	org := Organisation{
		Name:      "o_test_org2",
		MainGroup: "test2",
		SubGroup:  "test2",
		Flair:     sql.NullString{Valid: true, String: "test"},
		Status:    Public,
		Members:   []Account{acc3},
		Admins:    []Account{acc1},
		Accounts:  []Account{acc2},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = Organisation{
		Name:      "o_test_org3",
		MainGroup: "test2",
		SubGroup:  "test2",
		Flair:     sql.NullString{Valid: true, String: "test"},
		Status:    Secret,
		Members:   []Account{acc2},
		Admins:    []Account{acc3},
		Accounts:  []Account{acc1},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = Organisation{
		Name:      "o_test_org_hidden",
		MainGroup: "test",
		SubGroup:  "test",
		Status:    Hidden,
		Members:   []Account{},
		Admins:    []Account{},
		Accounts:  []Account{},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
}

func testGetByAccounts(t *testing.T) {
	var member, admin, account int64
	acc := Account{}
	err := acc.GetByUserName("test_org1")
	assert.Nil(t, err)
	member = acc.ID
	err = acc.GetByUserName("test_org2")
	assert.Nil(t, err)
	admin = acc.ID
	err = acc.GetByUserName("test_org3")
	assert.Nil(t, err)
	account = acc.ID
	org := Organisation{}
	err = org.GetByName("o_test_org")
	assert.Nil(t, err)
	//member
	err = org.GetByNameAndOnlyWhenAccountIsMember("o_test_org", member)
	assert.Nil(t, err)
	assert.Equal(t, "test", org.SubGroup)
	err = org.GetByNameAndOnlyWhenAccountIsMember("o_test_org", admin)
	assert.Nil(t, err)
	assert.Equal(t, "test", org.SubGroup)
	err = org.GetByNameAndOnlyWhenAccountIsMember("o_test_org", account)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	//admin
	err = org.GetByNameAndOnlyWhenAccountAsAdmin("o_test_org", member)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = org.GetByNameAndOnlyWhenAccountAsAdmin("o_test_org", admin)
	assert.Nil(t, err)
	assert.Equal(t, "test", org.SubGroup)
	err = org.GetByNameAndOnlyWhenAccountAsAdmin("o_test_org", account)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	//account
	err = org.GetByNameAndOnlyWithAccount("o_test_org", member)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = org.GetByNameAndOnlyWithAccount("o_test_org", admin)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = org.GetByNameAndOnlyWithAccount("o_test_org", account)
	assert.Nil(t, err)
	assert.Equal(t, "test", org.SubGroup)
}

func testOrganisationChange(t *testing.T) {
	org := Organisation{}
	err := org.GetByName("o_test_org")
	assert.Nil(t, err)
	org.MainGroup = "test2"
	org.Status = Private
	acc := Account{}
	err = acc.GetByUserName("test_org1")
	assert.Nil(t, err)
	org.Members = []Account{acc}
	err = acc.GetByUserName("test_org2")
	assert.Nil(t, err)
	org.Admins = []Account{acc}
	err = acc.GetByUserName("test_org3")
	assert.Nil(t, err)
	org.Accounts = []Account{acc}
	err = org.SaveChanges()
	assert.Nil(t, err)
	err = org.UpdateAdmins()
	assert.Nil(t, err)
	err = org.UpdateAccounts()
	assert.Nil(t, err)
	err = org.UpdateMembers()
	assert.Nil(t, err)
	second := Organisation{}
	err = second.GetByName("o_test_org")
	assert.Nil(t, err)
	assert.Equal(t, org, second)
}

func testOrgansationCreate(t *testing.T) {
	org := Organisation{
		Name:      "o_test_org",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{Valid: true, String: "test"},
		Status:    Public,
		Members:   []Account{},
		Admins:    []Account{},
		Accounts:  []Account{},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	org2 := Organisation{}
	err = org2.GetByName("o_test_org")
	assert.Nil(t, err)
	assert.Equal(t, org, org2)
}

func testSetupAccountsOrganisation(t *testing.T) {
	forTestCreateAccount(t, "test_org1", Account{Role: User})
	forTestCreateAccount(t, "test_org2", Account{Role: User})
	forTestCreateAccount(t, "test_org3", Account{Role: User})
}
