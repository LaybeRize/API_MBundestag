package htmlHandler

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	wr "API_MBundestag/htmlWrapper"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
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

func FillAllNotSuspendedNames[T any](v *T) {
	names, err := dataLogic.GetAllAccountNamesNotSuspended()
	slice := reflect.ValueOf(names)
	updateField(v, "Names", slice, err, generics.NamesQueryError)
}

func FillUserAndDisplayNames[T any](v *T) {
	names := database.NameList{}
	err := names.GetAllUserAndDisplayName()
	slice := reflect.ValueOf(names)
	updateField(v, "Names", slice, err, generics.NamesQueryError)
}

func FillOrganisationNames[T any](v *T) {
	orgNames, err := dataLogic.GetAllOrganisationNames()
	slice := reflect.ValueOf(orgNames)
	updateField(v, "OrgNames", slice, err, generics.OrgNamesQueryError)
}

func FillOrganisationGroups[T any](v *T) {
	main, sub, err := dataLogic.GetNamesForSubAndMainGroups()
	sliceMain := reflect.ValueOf(main)
	sliceSub := reflect.ValueOf(sub)
	updateField(v, "ExistingMainGroup", sliceMain, nil, "")
	updateField(v, "ExistingSubGroup", sliceSub, err, generics.GroupQueryError)
}

func FillOwnAccounts[T any](v *T, acc *database.Account) {
	ownAccounts := database.AccountList{}
	err := ownAccounts.GetAllPressAccountsFromAccountPlusSelf(acc)
	slice := reflect.ValueOf(ownAccounts)
	updateField(v, "Accounts", slice, err, generics.OwnAccountsCouldNotBeFound)
}

func FillOwnOrganisations[T any](v *T, acc *database.Account) {
	ownOrgs := database.OrganisationList{}
	var err error
	if acc.Role == database.HeadAdmin {
		err = ownOrgs.GetAllVisable()
	} else {
		err = ownOrgs.GetAllPartOf(acc.ID)
	}
	slice := reflect.ValueOf(ownOrgs)
	updateField(v, "Organisations", slice, err, generics.OrgNamesQueryError)
}

func updateField[T any](v *T, name string, slice reflect.Value, err error, errorText string) {
	ref := reflect.ValueOf(v).Elem()
	ref.FieldByName(name).Set(slice)
	if err != nil {
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(errorText + "\n" + mesg)
	}
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

func getFunctionName(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	return strs[len(strs)-1]
}

func AddFunctionToLinks(link string, function gin.HandlerFunc) {
	isPost := strings.HasPrefix(getFunctionName(function), "Post")
	Links = append(Links, Routing{
		IsPost: isPost,
		HFunc:  function,
		Link:   link,
	})
}

type Routing struct {
	IsPost bool
	HFunc  gin.HandlerFunc
	Link   string
}

var Links = make([]Routing, 0)
