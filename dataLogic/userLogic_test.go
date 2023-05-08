package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserLogic(t *testing.T) {
	database.TestSetup()

	t.Run("testTestSetupAccount", testTestSetupAccount)
	t.Run("testGetUser", testGetUser)
	t.Run("testChangeUser", testChangeUser)
	t.Run("testChangePassword", testChangePassword)
	t.Run("testChangeLoginTries", testChangeLoginTries)
	t.Run("testUpdateResetLoginTries", testUpdateResetLoginTries)
}

func testUpdateResetLoginTries(t *testing.T) {
	err := ResetLoginTries("accUserLogic")
	assert.Nil(t, err)
	acc := database.Account{}
	err = acc.GetByUserName("accUserLogic")
	assert.Nil(t, err)
	assert.Equal(t, 0, acc.LoginTries)
	assert.False(t, acc.NextLoginTime.Valid)
}

func testChangeLoginTries(t *testing.T) {
	acc := database.Account{Username: "accUserLogic"}
	err := UpdateLoginTries(&acc)
	assert.Nil(t, err)
	err = acc.GetByDisplayName(acc.DisplayName)
	assert.Nil(t, err)
	assert.Equal(t, 1, acc.LoginTries)
	assert.False(t, acc.NextLoginTime.Valid)
	err = UpdateLoginTries(&acc)
	assert.Nil(t, err)
	err = UpdateLoginTries(&acc)
	assert.Nil(t, err)
	err = UpdateLoginTries(&acc)
	assert.Equal(t, AccountCanNotBeLoggindBecauseOfTimeout, err)
	err = acc.GetByDisplayName(acc.DisplayName)
	assert.Nil(t, err)
	assert.Equal(t, 4, acc.LoginTries)
	assert.True(t, acc.NextLoginTime.Valid)
}

func testChangePassword(t *testing.T) {
	var msg generics.Message
	var positive bool
	ChangePassword("askjdlasd", "", "", &msg, &positive)
	assert.Equal(t, AccountCloudNotBeFound+"\n", msg)
	assert.False(t, positive)
	msg = ""

	ChangePassword("accUserLogic", "asdvasdasd", "", &msg, &positive)
	assert.Equal(t, OldPasswordNotcorrect+"\n", msg)
	assert.False(t, positive)
	msg = ""

	ChangePassword("accUserLogic", "test", "testNew", &msg, &positive)
	assert.Equal(t, AccountPasswordSuccessfulChanged+"\n", msg)
	assert.True(t, positive)

	accDB := database.Account{}
	err := accDB.GetByUserName("accUserLogic")
	assert.Nil(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(accDB.Password), []byte("testNew"))
	assert.Nil(t, err)
}

func testChangeUser(t *testing.T) {
	acc := Account{}
	var msg generics.Message
	var positive bool
	acc.GetUser("accUserLogic", "", &msg, &positive)
	acc.ChangeFlair = true
	acc.Flair = "testFlair set now"
	acc.Role = database.PressAccount
	acc.Linked = 1
	acc.RemoveFromTitle = true
	acc.RemoveFromOrganisation = true

	msg = ""
	positive = false
	acc.ChangeUser(&msg, &positive)
	assert.Equal(t, CouldChangeAccount+"\n", msg)
	assert.True(t, positive)

	accDB := database.Account{}
	err := accDB.GetByUserName("accUserLogic")
	assert.Nil(t, err)
	assert.Equal(t, "accUserLogic", accDB.DisplayName)
	assert.Equal(t, "accUserLogic", accDB.Username)
	assert.Equal(t, "testFlair set now", accDB.Flair)
	assert.Equal(t, database.PressAccount, accDB.Role)
	assert.True(t, accDB.Linked.Valid)
	assert.Equal(t, int64(1), accDB.Linked.Int64)
}

func testGetUser(t *testing.T) {
	acc := Account{}
	var msg generics.Message
	var positive bool
	acc.GetUser("accUserLogic", "a", &msg, &positive)
	assert.Equal(t, CouldNotFindAccount+"\n", msg)
	assert.False(t, positive)
	msg = ""
	acc.GetUser("234zhsgdfsf", "", &msg, &positive)
	assert.Equal(t, CouldNotFindAccount+"\n", msg)
	assert.False(t, positive)
	msg = ""

	acc.GetUser("accUserLogic", "", &msg, &positive)
	assert.Equal(t, CouldFindAccount+"\n", msg)
	assert.True(t, positive)
	acc.ID = 0
	assert.Equal(t, Account{
		DisplayName: "accUserLogic",
		Username:    "accUserLogic",
		Role:        database.User,
	}, acc)
	acc.GetUser("accUserLogic", "", &msg, &positive)

	acc2 := Account{}
	msg = ""
	positive = false
	acc2.GetUser("", "accUserLogic", &msg, &positive)
	assert.Equal(t, CouldFindAccount+"\n", msg)
	assert.True(t, positive)

	assert.Equal(t, acc, acc2)
}

func testTestSetupAccount(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "accUserLogic",
		Username:    "accUserLogic",
		Password:    string(hash),
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
}
