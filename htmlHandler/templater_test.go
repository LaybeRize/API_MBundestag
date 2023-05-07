package htmlHandler

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	wr "API_MBundestag/htmlWrapper"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestErrorValidator(t *testing.T) {
	temp := ValidationErrors{Info: "test"}
	assert.Equal(t, "test", temp.Error())
}

func TestMiddleHardwareForTests(t *testing.T) {
	ChangePath()
	w, ctx := GetEmptyContext(t)
	MiddleHardwareForTests(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestTemplating(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Title}} {{.Page.Info}}")
	assert.Nil(t, err)
	SetTemplate(t, temp, "testTemplating")

	t.Run("correctIdentity", correctIdentity)
	t.Run("correctTemplating", correctTemplating)
}

func correctTemplating(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	type TestStruct struct {
		LoggedIn bool
		Account  database.Account
		Info     string
	}
	PageIdentityMap[Identity(TestStruct{})] = BasicStruct{
		Title:    "TestSite",
		Site:     "testTemplating",
		Template: "testTemplating",
	}
	MakeSite(&TestStruct{Info: "adjknlödf"}, ctx, &database.Account{})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "TestSite adjknlödf", w.Body.String())
}

func correctIdentity(t *testing.T) {
	type TestStruct struct {
		LoggedIn bool
		Account  database.Account
		Info     string
	}
	assert.Equal(t, "htmlHandler.TestStruct", Identity(TestStruct{}))
}

func TestReflection(t *testing.T) {
	type TestStruct struct {
		TestField []int
		generics.MessageStruct
	}
	ts := TestStruct{}
	slice := reflect.ValueOf([]int{1, 2, 3})
	updateField(&ts, "TestField", slice, nil, "")
	assert.Equal(t, TestStruct{TestField: []int{1, 2, 3}}, ts)

	slice = reflect.ValueOf([]int{4, 5, 6})
	ts.Message = "Test\n"
	updateField(&ts, "TestField", slice, errors.New("example"), "Error")
	assert.Equal(t, TestStruct{TestField: []int{4, 5, 6}, MessageStruct: generics.MessageStruct{
		Message: "Error\nTest\n",
		Positiv: false,
	}}, ts)
}
