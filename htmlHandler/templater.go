package htmlHandler

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	wr "API_MBundestag/htmlWrapper"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

var Template *wr.Templates

type ValidationErrors struct {
	Info string
}

func (err ValidationErrors) Error() string {
	return err.Info
}

type BasicStruct struct {
	Title    string
	Site     string
	Template string
	User     database.Account
	Page     interface{}
}

var templateDir = "templates"

func MiddleHardwareForTests(c *gin.Context) {
	var err error
	Template, err = wr.New(templateDir, ".html", wr.DefaultFunctions)
	if err != nil {
		log.Fatal(err)
	}
	help.UpdateAttributes()
	c.Redirect(http.StatusTemporaryRedirect, c.Param("path"))
}

/* FOR SIMPLIFICATIONS */

func ExtractAmount(c *gin.Context, min int, max int, standard int) int {
	i, err := strconv.Atoi(c.Query("amount"))
	if err != nil {
		return standard
	}
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}

func FillTitleNames[T any](v *T) {
	ref := reflect.ValueOf(v).Elem()
	titleNames := reflect.ValueOf(dataLogic.GetTitleNames())
	ref.FieldByName("TitleNames").Set(titleNames)
}

func FillTitleGroups[T any](v *T) {
	ref := reflect.ValueOf(v).Elem()
	mainGroup := reflect.ValueOf(dataLogic.GetMainGroupNames())
	subGroup := reflect.ValueOf(dataLogic.GetSubGroupNames())
	ref.FieldByName("ExistingMainGroup").Set(mainGroup)
	ref.FieldByName("ExistingSubGroup").Set(subGroup)
}

func Identity[T any](v T) string {
	return fmt.Sprintf("%T", v)
}

func MakeSite[T any](v *T, c *gin.Context, acc *database.Account) {
	info := PageIdentityMap[Identity(*v)]
	info.User = *acc
	info.Page = v
	err := Template.Render(c.Writer, info.Template, info, http.StatusOK)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

var PageIdentityMap = map[string]BasicStruct{
	//Approved
	Identity(dataLogic.TitleMainGroupArray{}): {
		Title:    "Titlesübersicht",
		Site:     "displayTitle",
		Template: "displayTitle",
	},
	//Approved
	Identity(dataLogic.OrganisationMainGroupArray{}): {
		Title:    "Organisationsübersicht",
		Site:     "displayOrganisation",
		Template: "displayOrganisation",
	},
}
