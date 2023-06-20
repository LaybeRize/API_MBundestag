package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type CreateTitleStruct struct {
	Title             dataLogic.Title
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	generics.MessageStruct
}

func getEmptyCreateTitleStruct() *CreateTitleStruct {
	request := CreateTitleStruct{}
	htmlHandler.FillAllNotSuspendedNames(&request)
	htmlHandler.FillTitleGroups(&request)
	return &request
}

func GetCreateTitlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(getEmptyCreateTitleStruct(), c, &acc)
}

func PostCreateTitlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(validateCreateTitle(c), c, &acc)
}

func validateCreateTitle(c *gin.Context) (titleStruct *CreateTitleStruct) {
	titleStruct = getEmptyCreateTitleStruct()
	titleStruct.Title = dataLogic.Title{
		Name:      generics.GetText(c, "name"),
		MainGroup: generics.GetText(c, "mainGroup"),
		SubGroup:  generics.GetText(c, "subGroup"),
		Flair:     generics.GetText(c, "flair"),
		Holder:    generics.GetStringArray(c, "user"),
	}

	switch true {
	case titleStruct.Message.CheckOrgOrTitle(&titleStruct.Title):
	default:
		titleStruct.Title.CreateMe(&titleStruct.Message, &titleStruct.Positiv)
	}
	return
}
