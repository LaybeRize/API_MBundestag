package htmlHandler

import (
	"API_MBundestag/database"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func GetEmptyContext(t *testing.T) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{}
	ctx.Request.Header = http.Header{}
	ctx.Request.URL = &url.URL{}
	err := ctx.Request.ParseForm()
	assert.Nil(t, err)
	return w, ctx
}

func GetContextWithForm(t *testing.T, m map[string]string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{}
	ctx.Request.URL = &url.URL{}
	ctx.Request.Header = http.Header{}
	err := ctx.Request.ParseForm()
	assert.Nil(t, err)
	for key, value := range m {
		ctx.Request.PostForm.Add(key, value)
	}
	return w, ctx
}

func GetContextSetForUser(t *testing.T, acc database.Account) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{}
	ctx.Request.Header = http.Header{}
	ctx.Request.URL = &url.URL{}
	err := ctx.Request.ParseForm()
	assert.Nil(t, err)
	acc.RefToken.Valid = true
	acc.RefToken.String = "abc"
	acc.ExpDate.Valid = true
	acc.ExpDate.Time = time.Now().UTC().Add(time.Hour)
	ctx.Request.Header = http.Header{}
	ctx.Request.AddCookie(&http.Cookie{
		Name:    "token",
		Value:   "abc",
		Expires: time.Now().UTC().Add(time.Hour),
	})
	err = acc.SaveChanges()
	assert.Nil(t, err)
	return w, ctx
}

func GetContextSetForUserWithForm(t *testing.T, acc database.Account, m map[string]string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{}
	ctx.Request.Header = http.Header{}
	ctx.Request.URL = &url.URL{}
	err := ctx.Request.ParseForm()
	assert.Nil(t, err)
	for key, value := range m {
		ctx.Request.PostForm.Add(key, value)
	}
	assert.Nil(t, err)
	acc.RefToken.Valid = true
	acc.RefToken.String = uuid.New().String()
	acc.ExpDate.Valid = true
	acc.ExpDate.Time = time.Now().UTC().Add(time.Hour)
	ctx.Request.Header = http.Header{}
	ctx.Request.AddCookie(&http.Cookie{
		Name:    "token",
		Value:   acc.RefToken.String,
		Expires: time.Now().UTC().Add(time.Hour),
	})
	err = acc.SaveChanges()
	assert.Nil(t, err)
	return w, ctx
}

func CreateAccountForTest(t *testing.T, name string, role database.RoleString, link int64, flair ...string) {
	acc := database.Account{
		Username:    name,
		DisplayName: name,
		Role:        role,
		Linked:      sql.NullInt64{Valid: link != 0, Int64: link},
	}
	if len(flair) == 1 {
		acc.Flair = flair[0]
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
}

func PadLeft(s, p string, count int) string {
	ret := make([]byte, len(p)*count+len(s))

	b := ret[:len(p)*count]
	bp := copy(b, p)
	for bp < len(b) {
		copy(b[bp:], b[:bp])
		bp *= 2
	}
	copy(ret[len(b):], s)
	return string(ret)
}

var templateEditing = sync.Mutex{}

func SetTemplate(t *testing.T, temp *template.Template, name string) {
	templateEditing.Lock()
	defer templateEditing.Unlock()
	if Template == nil {
		errTemp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestError|{{.Page.Error}}")
		assert.Nil(t, err)
		Template = &wr.Templates{
			Extension: "",
			Dir:       "",
			Templates: map[string]*template.Template{
				"error": errTemp,
			},
		}
	}
	Template.Templates[name] = temp
}

func GetContextSetForUserWithFormAndQuery(t *testing.T, acc database.Account, m map[string]string, query string) (*httptest.ResponseRecorder, *gin.Context) {
	w, ctx := GetContextSetForUserWithForm(t, acc, m)
	ctx.Request.URL.RawQuery = query
	return w, ctx
}

func ChangePath() {
	_, filename, _, _ := runtime.Caller(0)
	pathStr := path.Dir(filename)
	if strings.HasSuffix(pathStr, "API_MBundestag") {
		return
	}
	var re = regexp.MustCompile(`(?m)^.*API_MBundestag`)

	pathStr = re.FindAllString(pathStr, -1)[0]

	err := os.Chdir(pathStr) // change to suit test file location
	if err != nil {
		log.Println(err)
	}
}
