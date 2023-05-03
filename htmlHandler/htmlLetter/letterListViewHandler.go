package htmlLetter

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type ViewLetterListStruct struct {
	Search          bool
	HasNext         bool
	HasBefore       bool
	NextUUID        string
	BeforeUUID      string
	Amount          int
	LetterList      database.LetterList
	SelectedAccount string
	Accounts        database.AccountList
	FormatString    string
	Message         string
}

type ViewModMailListStrcut struct {
	ViewLetterListStruct
}

func getEmtpyViewLetterListStruct(acc *database.Account) *ViewLetterListStruct {
	val := ViewLetterListStruct{}
	htmlHandler.FillOwnAccounts(&val, acc)
	val.Search = true
	val.FormatString = generics.LongTimeString
	return &val
}

func getEmtpyViewModMailListStruct() *ViewLetterListStruct {
	val := ViewLetterListStruct{}
	val.Search = false
	val.FormatString = generics.LongTimeString
	return &val
}

func GetViewModMailListPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	i := htmlHandler.ExtractAmount(c, 1, 50, 20)

	var err error
	letterStruct := getEmtpyViewModMailListStruct()
	if generics.GetIfType(c, "before") {
		err = letterStruct.validateLetterReadPageBefore(c, i, "", true)
	} else {
		err = letterStruct.validateLetterReadNextPage(c, i, "", true)
	}
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLoadingLetters)
		return
	}

	letterStruct.Amount = i
	htmlHandler.MakeSite(&ViewModMailListStrcut{*letterStruct}, c, &acc)
}

func GetViewLetterListPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	//check viewer
	viewer := database.Account{}
	err := viewer.GetByDisplayName(c.Query("usr"))
	if err != nil && !generics.GetIfEmptyQuery(c, "usr") {
		htmlBasics.MakeErrorPage(c, &acc, generics.AccountDoesNotExistOrIsNotYours)
		return
	}
	if generics.GetIfEmptyQuery(c, "usr") {
		viewer = acc
	}
	if viewer.DisplayName != acc.DisplayName && viewer.Linked.Int64 != acc.ID {
		htmlBasics.MakeErrorPage(c, &acc, generics.AccountDoesNotExistOrIsNotYours)
		return
	}

	i := htmlHandler.ExtractAmount(c, 1, 50, 20)

	letterStruct := getEmtpyViewLetterListStruct(&acc)
	letterStruct.SelectedAccount = viewer.DisplayName
	if generics.GetIfType(c, "before") {
		err = letterStruct.validateLetterReadPageBefore(c, i, viewer.DisplayName, false)
	} else {
		err = letterStruct.validateLetterReadNextPage(c, i, viewer.DisplayName, false)
	}
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLoadingLetters)
		return
	}

	letterStruct.Amount = i
	htmlHandler.MakeSite(letterStruct, c, &acc)
}

func PostViewLetterListPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	c.Redirect(http.StatusFound, "/letter-list?usr="+url.QueryEscape(c.PostForm("selectedAccount")))
}

func (listStruct *ViewLetterListStruct) validateLetterReadNextPage(c *gin.Context, i int, acc string, modMails bool) error {
	err, exists := listStruct.LetterList.GetPublicationAfter(c.Query("uuid"), i+1, acc, modMails)
	if len(listStruct.LetterList) == 0 {
		return err
	}
	if len(listStruct.LetterList) == i+1 {
		listStruct.HasNext = true
		listStruct.NextUUID = listStruct.LetterList[i-1].UUID
		listStruct.LetterList = listStruct.LetterList[:i]
	}
	if exists {
		listStruct.HasBefore = true
		listStruct.BeforeUUID = listStruct.LetterList[0].UUID
	}
	return err
}

func (listStruct *ViewLetterListStruct) validateLetterReadPageBefore(c *gin.Context, i int, acc string, modMails bool) error {
	err, exists := listStruct.LetterList.GetPublicationBefore(c.Query("uuid"), i+1, acc, modMails)
	if len(listStruct.LetterList) == 0 {
		return err
	}
	if len(listStruct.LetterList) == i+1 {
		listStruct.HasBefore = true
		listStruct.LetterList = listStruct.LetterList[1:]
		listStruct.BeforeUUID = listStruct.LetterList[0].UUID
	}
	if exists {
		listStruct.HasNext = true
		listStruct.NextUUID = listStruct.LetterList[len(listStruct.LetterList)-1].UUID
	}
	return err
}
