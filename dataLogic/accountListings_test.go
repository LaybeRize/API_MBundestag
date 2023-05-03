package dataLogic

import (
	"API_MBundestag/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountListings(t *testing.T) {
	database.TestSetup()

	acc := database.Account{
		DisplayName: "accInfo_a",
		Username:    "accInfo_a",
		Password:    "test",
		Role:        database.User,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.ID = 0
	acc.DisplayName, acc.Username = "accInfo_b", "accInfo_b"
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.ID = 0
	acc.DisplayName, acc.Username = "accInfo_c", "accInfo_c"
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.ID = 0
	acc.DisplayName, acc.Username = "accInfo_d", "accInfo_d"
	err = acc.CreateMe()
	assert.Nil(t, err)
	array, err := GetAllAccountNamesNotSuspended()
	assert.Nil(t, err)
	counter := 0
	for _, name := range array {
		switch name {
		case "accInfo_a", "accInfo_b", "accInfo_c", "accInfo_d":
			counter++
		}
	}
	assert.Equal(t, 4, counter)
	acc.Suspended = true
	err = acc.SaveChanges()
	assert.Nil(t, err)
	array, err = GetAllAccountNamesNotSuspended()
	assert.Nil(t, err)
	counter = 0
	for _, name := range array {
		switch name {
		case "accInfo_a", "accInfo_b", "accInfo_c", "accInfo_d":
			counter++
		}
	}
	assert.Equal(t, 3, counter)
}
