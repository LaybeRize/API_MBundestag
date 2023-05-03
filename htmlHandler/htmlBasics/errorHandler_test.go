package htmlBasics

import (
	"API_MBundestag/database"
	"API_MBundestag/htmlHandler"
	hHa "API_MBundestag/htmlHandler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	database.TestSetup()
	Setup()

	acc := database.Account{
		DisplayName: "test_errorHandlingHTML",
		Username:    "test_errorHandlingHTML",
		Password:    "test",
		Role:        database.User,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	hHa.SetTemplate(t, nil, "")

	w, ctx := htmlHandler.GetEmptyContext(t)
	MakeErrorPage(ctx, &acc, "test")

	assert.Equal(t, "TestError|test", w.Body.String())
}
