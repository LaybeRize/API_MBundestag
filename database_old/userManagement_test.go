package database

import (
	"API_MBundestag/help"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expected = Account{
	ID:          1,
	DisplayName: "test",
	Flair:       "tevcyxc",
	Username:    "bsdfsad",
	Password:    "tsdfsdaf",
	Suspended:   false,
	RefToken:    sql.NullString{},
	ExpDate:     sql.NullTime{},
	Role:        "user",
	Linked: sql.NullInt64{
		Valid: true,
		Int64: 12,
	},
}

func TestAccount(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestAccountDB()

	t.Run("testCreateAccount", testCreateAccount)
	t.Run("testGetAccountByDisplayname", testGetAccountByDisplayname)
	t.Run("testGetAccountByUsername", testGetAccountByUsername)
	t.Run("testEditAccount", testEditAccount)
	t.Run("testGetAccountByRefreshToken", testGetAccountByRefreshToken)
	t.Run("testGetAccountByID", testGetAccountByID)
	t.Run("testCheckIfAccountsExists", testCheckIfAccountsExists)
	t.Run("testCheckPressAccounts", testCheckPressAccounts)
	t.Run("testIfAccountGetsOverwritenOnSQLRowNotFound", testIfAccountGetsOverwritenOnSQLRowNotFound)
}

func testIfAccountGetsOverwritenOnSQLRowNotFound(t *testing.T) {
	acc := Account{
		ID:          54323,
		DisplayName: "asd",
		Flair:       "ycx",
		Username:    "xcvse",
		Password:    "dsfad",
		Suspended:   true,
		RefToken:    sql.NullString{},
		ExpDate:     sql.NullTime{},
		Role:        "",
		Linked: sql.NullInt64{
			Valid: true,
			Int64: 12,
		},
	}
	err := acc.GetByUserName("sjhbdkhoruewir")
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Equal(t, int64(54323), acc.ID)
	assert.Equal(t, "xcvse", acc.Username)
	assert.Equal(t, "dsfad", acc.Password)
}

func testCheckPressAccounts(t *testing.T) {
	list := AccountList{}
	acc := Account{}
	err := acc.GetByDisplayName("other")
	assert.Nil(t, err)
	acc.Suspended = false
	err = acc.SaveChanges()
	assert.Nil(t, err)
	err = acc.GetByDisplayName("test")
	err = list.GetAllPressAccountsFromAccountPlusSelf(acc)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test", list[0].DisplayName)
	assert.Equal(t, "other", list[1].DisplayName)
}

func testCheckIfAccountsExists(t *testing.T) {
	b, err := DoAccountsExist([]string{"test", "other"}, false)
	assert.Equal(t, false, b)
	assert.Equal(t, errors.New("other"), err)
	b, err = DoAccountsExist([]string{"test", "other", "third"}, false)
	assert.Equal(t, false, b)
	assert.Equal(t, errors.New("other"), err)
	acc := Account{
		DisplayName: "other",
		Username:    "other",
		Password:    "test",
		Role:        PressAccount,
		Linked: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.Suspended = true
	err = acc.SaveChanges()
	assert.Nil(t, err)
	b, err = DoAccountsExist([]string{"test", "other"}, false)
	assert.Equal(t, errors.New("other"), err)
	assert.False(t, b)
	b, err = DoAccountsExist([]string{"test", "other"}, true)
	assert.Nil(t, err)
	assert.True(t, b)
}

func testGetAccountByID(t *testing.T) {
	expected = Account{
		ID:          1,
		DisplayName: "test",
		Flair:       "tevcyxc",
		Username:    "bsdfsad",
		Password:    "tsdfsdaf",
		Suspended:   false,
		RefToken: sql.NullString{
			String: "testRefToken",
			Valid:  true,
		},
		ExpDate: sql.NullTime{},
		Role:    "user",
		Linked: sql.NullInt64{
			Valid: true,
			Int64: 12,
		},
	}
	acc := Account{}
	err := acc.GetByID(1)
	assert.Nil(t, err)

	assert.Equal(t, expected, acc)
}

func testGetAccountByRefreshToken(t *testing.T) {
	expected = Account{
		ID:          1,
		DisplayName: "test",
		Flair:       "tevcyxc",
		Username:    "bsdfsad",
		Password:    "tsdfsdaf",
		Suspended:   false,
		RefToken: sql.NullString{
			String: "testRefToken",
			Valid:  true,
		},
		ExpDate: sql.NullTime{},
		Role:    "user",
		Linked: sql.NullInt64{
			Valid: true,
			Int64: 12,
		},
	}
	acc := Account{}
	err := acc.GetByToken("testRefToken")
	assert.Nil(t, err)

	assert.Equal(t, expected, acc)
}

func testEditAccount(t *testing.T) {
	acc := Account{}
	err := acc.GetByUserName("bsdfsad")
	assert.Nil(t, err)

	acc.RefToken.String = "testRefToken"
	acc.RefToken.Valid = true
	err = acc.SaveChanges()
	assert.Nil(t, err)

	result := Account{}
	err = result.GetByUserName("bsdfsad")
	assert.Nil(t, err)

	assert.Equal(t, acc, result)
}

func testGetAccountByDisplayname(t *testing.T) {
	acc := Account{}
	err := acc.GetByDisplayName("test")
	assert.Nil(t, err)

	assert.Equal(t, expected, acc)
}

func testGetAccountByUsername(t *testing.T) {
	acc := Account{}
	err := acc.GetByUserName("bsdfsad")
	assert.Nil(t, err)

	assert.Equal(t, expected, acc)
}

func testCreateAccount(t *testing.T) {
	acc := Account{
		ID:          3,
		DisplayName: "test",
		Flair:       "tevcyxc",
		Username:    "bsdfsad",
		Password:    "tsdfsdaf",
		Suspended:   true,
		RefToken: sql.NullString{
			Valid:  true,
			String: "test",
		},
		ExpDate: sql.NullTime{},
		Role:    "user",
		Linked: sql.NullInt64{
			Valid: true,
			Int64: 12,
		},
	}
	err := acc.CreateMe()
	assert.Nil(t, err)

	request := Account{}
	err = request.GetByUserName("bsdfsad")
	assert.Nil(t, err)

	assert.NotEqual(t, acc, request)

	assert.Equal(t, expected, request)
}
