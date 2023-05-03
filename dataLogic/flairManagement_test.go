package dataLogic

import (
	"API_MBundestag/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestFlairManager(t *testing.T) {
	database.TestSetup()

	t.Run("setUpAccountsTestFlair", setUpAccountsTestFlair)
	t.Run("testAddFlair", testAddFlair)
	t.Run("testRemoveFlair", testRemoveFlair)
}

var testFlairAccount *database.Account

func testAddFlair(t *testing.T) {
	query := database.Account{}
	err := addFlair("test", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test", testFlairAccount.Flair)

	err = query.GetByID(testFlairAccount.ID)
	assert.Nil(t, err)
	assert.Equal(t, "test", query.Flair)

	err = addFlair("test2", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test, test2", testFlairAccount.Flair)
	err = addFlair("test3", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test, test2, test3", testFlairAccount.Flair)
	err = addFlair("test4", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test, test2, test3, test4", testFlairAccount.Flair)

	err = query.GetByID(testFlairAccount.ID)
	assert.Nil(t, err)
	assert.Equal(t, "test, test2, test3, test4", query.Flair)
}

func testRemoveFlair(t *testing.T) {
	query := database.Account{}
	err := removeFlairWithSave("test2", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test, test3, test4", testFlairAccount.Flair)

	err = query.GetByID(testFlairAccount.ID)
	assert.Nil(t, err)
	assert.Equal(t, "test, test3, test4", query.Flair)

	err = removeFlairWithSave("test", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test3, test4", testFlairAccount.Flair)
	err = removeFlairWithSave("test4", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "test3", testFlairAccount.Flair)
	err = removeFlairWithSave("test3", testFlairAccount)
	assert.Nil(t, err)
	assert.Equal(t, "", testFlairAccount.Flair)

	err = query.GetByID(testFlairAccount.ID)
	assert.Nil(t, err)
	assert.Equal(t, "", query.Flair)
}

func setUpAccountsTestFlair(t *testing.T) {
	testFlairAccount = &database.Account{
		DisplayName: "test_flairManagemenet",
		Flair:       "",
		Username:    "test_flairManagemenet",
		Suspended:   false,
		Role:        database.User,
	}
	err := testFlairAccount.CreateMe()
	assert.Nil(t, err)
}
