package dataLogic

import (
	"API_MBundestag/database"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var orgHierarch = OrganisationMainGroupArray{}
var mainGroup1 = OrganisationMainGroup{}
var mainGroup2 = OrganisationMainGroup{}
var accHierarchy = database.Account{}

func TestOrganisationHierarchy(t *testing.T) {
	database.TestSetup()

	accHierarchy = database.Account{
		DisplayName: "orgHierarchy",
		Username:    "orgHierarchy",
		Password:    "orgHierarchy",
		Role:        database.HeadAdmin,
	}
	err := accHierarchy.CreateMe()
	assert.Nil(t, err)

	t.Run("testSingleOrgHierarchy", testSingleOrgHierarchy)
	t.Run("testSingleOrgInDifferentMainGroups", testSingleOrgInDifferentMainGroups)
	t.Run("testMultipleOrgInDifferentMainGroups", testMultipleOrgInDifferentMainGroups)
	t.Run("testMultipleOrgInDifferentMainGroupsWithSecret", testMultipleOrgInDifferentMainGroupsWithSecret)
	t.Run("testCorrectOptionHidden", testCorrectOptionHidden)
}

func testCorrectOptionHidden(t *testing.T) {
	org := database.Organisation{
		Name:      "test_hierarchy5",
		MainGroup: "test_hierarchy",
		SubGroup:  "test_hierarchy",
		Flair:     sql.NullString{},
		Status:    database.Hidden,
		Members:   []database.Account{},
		Admins:    []database.Account{},
		Accounts:  []database.Account{},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	mainGroup1.Amount = 1
	mainGroup1.Groups[0].Organisations = []database.Organisation{org}
	accHierarchy.Role = database.HeadAdmin
	err = orgHierarch.GetOrganisationHierarchy(accHierarchy, true)
	assert.Nil(t, err)
	orgHierarch.SetAmountForMainGroup()

	var ref *OrganisationMainGroup
	for _, m := range orgHierarch {
		if m.Name == "test_hierarchy" {
			ref = &m
			break
		}
	}
	assert.NotNil(t, ref)
	assert.Equal(t, &mainGroup1, ref)
}

func testMultipleOrgInDifferentMainGroupsWithSecret(t *testing.T) {
	org := database.Organisation{
		Name:      "test_hierarchy4",
		MainGroup: "test_hierarchy2",
		SubGroup:  "test_hierarchy2",
		Status:    database.Secret,
		Members:   []database.Account{accHierarchy},
		Admins:    []database.Account{},
		Accounts:  []database.Account{accHierarchy},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	mainGroup2.Groups[1].Organisations = append(mainGroup2.Groups[1].Organisations, org)
	mainGroup2.Groups[1].Amount = 2
	accHierarchy.Role = database.User
	err = orgHierarch.GetOrganisationHierarchy(accHierarchy, false)
	assert.Nil(t, err)

	var ref *OrganisationMainGroup
	var ref2 *OrganisationMainGroup
	for i, m := range orgHierarch {
		if m.Name == "test_hierarchy" {
			ref = &m
			ref2 = &orgHierarch[i+1]
			break
		}
	}
	assert.NotNil(t, ref)
	assert.Equal(t, &mainGroup1, ref)
	assert.NotNil(t, ref2)
	assert.Equal(t, &mainGroup2, ref2)
}

func testMultipleOrgInDifferentMainGroups(t *testing.T) {
	org := database.Organisation{
		Name:      "test_hierarchy3",
		MainGroup: "test_hierarchy2",
		SubGroup:  "test_hierarchy2",
		Status:    database.Public,
		Members:   []database.Account{},
		Admins:    []database.Account{},
		Accounts:  []database.Account{},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	mainGroup2.Groups = append(mainGroup2.Groups, OrganisationSubGroup{
		Name:          "test_hierarchy2",
		Amount:        1,
		Organisations: []database.Organisation{org},
	})
	err = orgHierarch.GetOrganisationHierarchy(accHierarchy, false)
	assert.Nil(t, err)

	var ref *OrganisationMainGroup
	var ref2 *OrganisationMainGroup
	for i, m := range orgHierarch {
		if m.Name == "test_hierarchy" {
			ref = &m
			ref2 = &orgHierarch[i+1]
			break
		}
	}
	assert.NotNil(t, ref)
	assert.Equal(t, &mainGroup1, ref)
	assert.NotNil(t, ref2)
	assert.Equal(t, &mainGroup2, ref2)
}

func testSingleOrgInDifferentMainGroups(t *testing.T) {
	org := database.Organisation{
		Name:      "test_hierarchy2",
		MainGroup: "test_hierarchy2",
		SubGroup:  "test_hierarchy",
		Status:    database.Public,
		Members:   []database.Account{},
		Admins:    []database.Account{},
		Accounts:  []database.Account{},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	mainGroup2 = OrganisationMainGroup{
		Name:   "test_hierarchy2",
		Amount: 0,
		Groups: []OrganisationSubGroup{
			{
				Name:          "test_hierarchy",
				Amount:        1,
				Organisations: []database.Organisation{org},
			},
		},
	}
	err = orgHierarch.GetOrganisationHierarchy(accHierarchy, false)
	assert.Nil(t, err)

	var ref *OrganisationMainGroup
	var ref2 *OrganisationMainGroup
	for i, m := range orgHierarch {
		if m.Name == "test_hierarchy" {
			ref = &m
			ref2 = &orgHierarch[i+1]
			break
		}
	}
	assert.NotNil(t, ref)
	assert.Equal(t, &mainGroup1, ref)
	assert.NotNil(t, ref2)
	assert.Equal(t, &mainGroup2, ref2)
}

func testSingleOrgHierarchy(t *testing.T) {
	org := database.Organisation{
		Name:      "test_hierarchy",
		MainGroup: "test_hierarchy",
		SubGroup:  "test_hierarchy",
		Status:    database.Public,
		Members:   []database.Account{},
		Admins:    []database.Account{},
		Accounts:  []database.Account{},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	mainGroup1 = OrganisationMainGroup{
		Name:   "test_hierarchy",
		Amount: 0,
		Groups: []OrganisationSubGroup{
			{
				Name:          "test_hierarchy",
				Amount:        1,
				Organisations: []database.Organisation{org},
			},
		},
	}
	err = orgHierarch.GetOrganisationHierarchy(accHierarchy, false)
	assert.Nil(t, err)

	var ref *OrganisationMainGroup
	for _, m := range orgHierarch {
		if m.Name == "test_hierarchy" {
			ref = &m
			break
		}
	}
	assert.NotNil(t, ref)
	assert.Equal(t, &mainGroup1, ref)
}
