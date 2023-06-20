package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type EditTitleStruct struct {
	TitleNames        []string
	Title             dataLogic.Title
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	generics.MessageStruct
}

func getEmptyEditTitleStruct() *EditTitleStruct {
	request := EditTitleStruct{}
	htmlHandler.FillAllNotSuspendedNames(&request)
	htmlHandler.FillTitleGroups(&request)
	htmlHandler.FillTitleNames(&request)
	return &request
}

func GetEditTitlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	titleStruct := getEmptyEditTitleStruct()
	if !generics.GetIfEmptyQuery(c, "title") {
		titleStruct.Title.GetMe(c.Query("title"))
	}

	htmlHandler.MakeSite(titleStruct, c, &acc)
}

func PostEditTitlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)

	switch true {
	case !b:
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
	case generics.GetIfType(c, "search"):
		htmlHandler.MakeSite(validateSearchTitle(c), c, &acc)
	case generics.GetIfType(c, "delete"):
		htmlHandler.MakeSite(validateDeleteTitle(c), c, &acc)
	default:
		htmlHandler.MakeSite(validateEditTitle(c), c, &acc)
	}
}

func validateSearchTitle(c *gin.Context) (editStruct *EditTitleStruct) {
	editStruct = getEmptyEditTitleStruct()
	err := editStruct.Title.GetMe(generics.GetText(c, "name"))
	if err != nil {
		editStruct.Message = generics.TitleDoesNotExists + "\n" + editStruct.Message
	} else {
		editStruct.Message = generics.SuccessfulFoundTitle + "\n" + editStruct.Message
		editStruct.Positiv = true
	}
	return
}

func validateDeleteTitle(c *gin.Context) (editStruct *EditTitleStruct) {
	editStruct = getEmptyEditTitleStruct()
	editStruct.Title.OldName = generics.GetText(c, "name")
	editStruct.Title.DeleteMe(&editStruct.Message, &editStruct.Positiv)
	return
}

func validateEditTitle(c *gin.Context) (editStruct *EditTitleStruct) {
	editStruct = getEmptyEditTitleStruct()
	editStruct.Title = dataLogic.Title{
		OldName:   generics.GetText(c, "name"),
		Name:      generics.GetText(c, "newName"),
		MainGroup: generics.GetText(c, "mainGroup"),
		SubGroup:  generics.GetText(c, "subGroup"),
		Flair:     generics.GetText(c, "flair"),
		Holder:    generics.GetStringArray(c, "user"),
	}

	switch true {
	case editStruct.Title.OldName == "":
		editStruct.Message = generics.NoNameForTitleProvided + "\n" + editStruct.Message
	case editStruct.Message.CheckOrgOrTitle(&editStruct.Title):
	default:
		editStruct.Title.ChangeMe(&editStruct.Message, &editStruct.Positiv)
	}
	return
}
