package dataLogic

import (
	"API_MBundestag/database"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrganisationNameList(t *testing.T) {
	database.TestSetup()

	t.Run("testCreateOrgs", testCreateOrgs)
	t.Run("testMainNamesOrgs", testNamesOrgs)
	t.Run("testSubNamesOrgs", testMainAndSubNamesOrgs)
}

func testMainAndSubNamesOrgs(t *testing.T) {
	main, sub, err := GetNamesForSubAndMainGroups()
	assert.Nil(t, err)
	counter := 0
	for _, m := range main {
		switch m {
		case "orgList_a", "orgList_b", "orgList_c", "orgList_d":
			counter++
		}
	}
	assert.Equal(t, 4, counter)
	counter = 0
	for _, m := range sub {
		switch m {
		case "orgList_a", "orgList_b", "orgList_c", "orgList_d":
			counter++
		}
	}
	assert.Equal(t, 4, counter)
}

func testNamesOrgs(t *testing.T) {
	names, err := GetAllOrganisationNames()
	assert.Nil(t, err)
	counter := 0
	for _, m := range names {
		switch m {
		case "orgList_a", "orgList_b", "orgList_c", "orgList_d":
			counter++
		}
	}
	assert.Equal(t, 4, counter)
}

func testCreateOrgs(t *testing.T) {
	org := database.Organisation{
		Name:      "orgList_a",
		MainGroup: "orgList_a",
		SubGroup:  "orgList_a",
		Flair:     sql.NullString{},
		Status:    database.Hidden,
		Members:   []database.Account{},
		Admins:    []database.Account{},
		Accounts:  []database.Account{},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.MainGroup, org.SubGroup = "orgList_b", "orgList_b", "orgList_b"
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.MainGroup, org.SubGroup = "orgList_c", "orgList_c", "orgList_c"
	org.Status = database.Public
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.MainGroup, org.SubGroup = "orgList_d", "orgList_d", "orgList_d"
	err = org.CreateMe()
	assert.Nil(t, err)
}
