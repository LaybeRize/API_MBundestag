package database

import (
	"API_MBundestag/help"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var org Organisation
var public OrganisationList
var private OrganisationList
var hidden OrganisationList

func TestOrganisations(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestOrganisationsDB()

	t.Run("testCreateorg", testCreateorg)
	t.Run("testGetOrgByName", testGetOrgByName)
	t.Run("testEditOrg", testEditOrg)
	t.Run("createOrgs", createOrgs)
	t.Run("testGetOrgsForNormalUser", testGetOrgsForNormalUser)
	t.Run("testGetOrgsForAdmins", testGetOrgsForAdmins)
	t.Run("testInvisableOrgs", testInvisableOrgs)
	t.Run("testMainGroupList", testMainGroupList)
	t.Run("testSubGroupList", testSubGroupList)
	t.Run("testFlairUniqueness", testFlairUniqueness)
	t.Run("testGetAllPartOf", testGetAllPartOf)
}

func testGetAllPartOf(t *testing.T) {
	orgList := OrganisationList{}
	err := orgList.GetAllPartOf("bazinga")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(orgList))
	err = orgList.GetAllPartOf("bazing")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(orgList))
	assert.Equal(t, "a", orgList[0].Name)
	assert.Equal(t, "lol", orgList[1].Name)
}

func testFlairUniqueness(t *testing.T) {
	org.Name = uuid.New().String()
	org.Flair = sql.NullString{
		String: "",
		Valid:  false,
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	org.Name = uuid.New().String()
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name = uuid.New().String()
	org.Flair = sql.NullString{
		String: "test",
		Valid:  true,
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name = uuid.New().String()
	err = org.CreateMe()
	assert.Equal(t, "pq: duplicate key value violates unique constraint \"organisations_flair_key\"", err.Error())
}

func testSubGroupList(t *testing.T) {
	orgList := OrganisationList{}
	err := orgList.GetAllSubGroups()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(orgList))
	assert.Equal(t, "dvasd", orgList[0].SubGroup)
}

func testMainGroupList(t *testing.T) {
	orgList := OrganisationList{}
	err := orgList.GetAllMainGroups()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(orgList))
	assert.Equal(t, "a", orgList[0].MainGroup)
}

func testInvisableOrgs(t *testing.T) {
	orgList := OrganisationList{}
	err := orgList.GetAllInvisable()
	assert.Nil(t, err)
	assert.Equal(t, hidden, orgList)
}

func testGetOrgsForAdmins(t *testing.T) {
	orgList := OrganisationList{}
	err := orgList.GetAllVisable()
	assert.Nil(t, err)
	assert.Equal(t, private, orgList)
}

func testGetOrgsForNormalUser(t *testing.T) {
	orgList := OrganisationList{}
	err := orgList.GetAllVisibleFor("")
	assert.Nil(t, err)
	assert.Equal(t, public, orgList)
}

func createOrgs(t *testing.T) {
	private = append(private, org)
	org.Flair = sql.NullString{}
	org.Name = "lol"
	org.Status = Public
	err := org.CreateMe()
	assert.Nil(t, err)
	private = append(private, org)
	public = append(public, org)
	org.Name = "zabruh"
	org.Status = Private
	org.Info.User = []string{}
	err = org.CreateMe()
	assert.Nil(t, err)
	private = append(private, org)
	public = append(public, org)
	org.Name = "bhda"
	org.Status = Hidden
	err = org.CreateMe()
	assert.Nil(t, err)
	hidden = append(hidden, org)
}

func testEditOrg(t *testing.T) {
	res := Organisation{}
	org.Flair.String = "asvxcv"
	org.Flair.Valid = true
	org.SubGroup = "dvasd"
	org.Status = Secret
	err := org.SaveChanges()
	assert.Nil(t, err)
	err = res.GetByName("a")
	assert.Nil(t, err)
	assert.Equal(t, org, res)
}

func testGetOrgByName(t *testing.T) {
	res := Organisation{}
	err := res.GetByName("a")
	assert.Nil(t, err)
	assert.Equal(t, org, res)
}

func testCreateorg(t *testing.T) {
	org = Organisation{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Status:    Public,
		Info: OrganisationInfo{
			Admins: []string{},
			User:   []string{"bazing", "ga"},
		},
	}
	err := org.CreateMe()
	assert.Nil(t, err)
}
