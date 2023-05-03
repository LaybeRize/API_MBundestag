package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var queryAccount = &Account{}

func TestUserManagment(t *testing.T) {
	TestSetup()
	t.Run("testCreateAllTypes", testCreateAllTypes)
	t.Run("testAllGetRequests", testAllGetRequests)
	t.Run("testSpecialQuerys", testSpecialQuerys)
	t.Run("testAccountChange", testAccountChange)
	t.Run("testLists", testLists)
}

func testCreateAllTypes(t *testing.T) {
	acc := &Account{
		DisplayName: "u_admin",
		Username:    "u_admin",
		Password:    "u_admin",
		Role:        Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc = &Account{
		DisplayName: "u_media_admin",
		Username:    "u_media_admin",
		Password:    "u_media_admin",
		Role:        MediaAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = &Account{
		DisplayName: "u_user",
		Username:    "u_user",
		Password:    "u_user",
		RefToken:    sql.NullString{Valid: true, String: "test"},
		Role:        User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = &Account{
		DisplayName: "u_press",
		Username:    "u_press",
		Password:    "u_press",
		Role:        PressAccount,
		Linked:      sql.NullInt64{Valid: true, Int64: 1},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
}

func testAllGetRequests(t *testing.T) {
	err := queryAccount.GetByDisplayName("u_user")
	assert.Nil(t, err)
	assert.Equal(t, "u_user", queryAccount.DisplayName)
	acc := &Account{}
	err = acc.GetByUserName("u_user")
	assert.Nil(t, err)
	assert.Equal(t, queryAccount, acc)
	err = acc.GetByDisplayName("u_user")
	assert.Nil(t, err)
	assert.Equal(t, queryAccount, acc)
	err = acc.GetByToken("test")
	assert.Nil(t, err)
	assert.Equal(t, queryAccount, acc)
	err = acc.GetByID(queryAccount.ID)
	assert.Nil(t, err)
	assert.Equal(t, queryAccount, acc)

}

func testSpecialQuerys(t *testing.T) {
	//check if I actually can aquire the parent
	acc := &Account{}
	err := acc.GetByDisplayNameWithParent("u_press")
	assert.Nil(t, err)
	assert.NotNil(t, acc.Parent)
	assert.Equal(t, "head_admin", acc.Parent.DisplayName)
	//check if the children function works too
	err = acc.GetByIDWithChildren(1)
	assert.Nil(t, err)
	assert.Equal(t, "head_admin", acc.DisplayName)
	assert.NotEqual(t, 0, len(acc.Children))
	var exists = false
	for _, item := range acc.Children {
		if item.DisplayName == "u_press" {
			exists = true
		}
	}
	assert.True(t, exists)

	list := AccountList{}
	err = list.GetAllPressAccountsFromAccountPlusSelf(&Account{ID: 1})
	assert.Nil(t, err)
	assert.True(t, len(list) > 1)
	assert.Equal(t, Account{ID: 1}, list[0])
	exists = false
	for _, item := range list {
		if item.DisplayName == "u_press" {
			exists = true
		}
	}
	assert.True(t, exists)
}

func testAccountChange(t *testing.T) {
	acc := &Account{}
	err := acc.GetByUserName("u_press")
	assert.Nil(t, err)
	assert.Equal(t, "u_press", acc.DisplayName)
	acc.Suspended = true
	err = acc.SaveChanges()
	assert.Nil(t, err)
	err = acc.GetByUserName("u_press")
	assert.Nil(t, err)
	assert.Equal(t, true, acc.Suspended)

	err = acc.GetByIDWithChildren(1)
	assert.Nil(t, err)
	assert.Equal(t, "head_admin", acc.DisplayName)
	var exists = false
	for _, item := range acc.Children {
		if item.DisplayName == "u_press" {
			exists = true
		}
	}
	assert.False(t, exists)
}

func testLists(t *testing.T) {
	list := AccountList{}
	names := NameList{}
	err := list.GetAllAccounts()
	assert.Nil(t, err)
	err = names.GetAllUserAndDisplayName()
	assert.Nil(t, err)
	counter := 0
	for i, acc := range list {
		assert.Equal(t, acc.DisplayName, names[i].DisplayName)
		switch acc.DisplayName {
		case "u_press":
			fallthrough
		case "u_user":
			fallthrough
		case "u_media_admin":
			fallthrough
		case "u_admin":
			fallthrough
		case "head_admin":
			counter++
		}
	}
	assert.Equal(t, 5, counter)

	err = list.GetAllAccountsNotSuspended()
	for _, acc := range list {
		if acc.DisplayName == "u_press" {
			t.Fail()
		}
	}

	var exists bool
	exists, err = list.DoAccountsExist([]string{"u_press", "u_user"})
	assert.False(t, exists)
	assert.Equal(t, "u_press", err.Error())
	assert.Equal(t, 0, len(list))
	exists, err = list.DoAccountsExist([]string{"u_media_admin", "u_user"})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "u_media_admin", list[0].DisplayName)
	assert.Equal(t, "u_user", list[1].DisplayName)
}

func forTestCreateAccount(t *testing.T, name string, account Account) {
	account.DisplayName = name
	account.Username = name
	account.Password = name
	err := account.CreateMe()
	assert.Nil(t, err)
}
