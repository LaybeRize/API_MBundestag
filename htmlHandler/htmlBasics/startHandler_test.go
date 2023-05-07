package htmlBasics

import (
	"API_MBundestag/database"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestStartPage(t *testing.T) {
	Setup()
	database.TestSetup()

	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "test",
		Username:    "test",
		Password:    string(hash),
		Role:        database.User,
		Linked:      sql.NullInt64{},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
}
