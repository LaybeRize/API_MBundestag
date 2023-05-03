package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	gen "API_MBundestag/htmlHandler/generics"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type CreateTitleStruct struct {
	Title             database.Title
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	Message           string
}

func getEmptyCreateTitleStruct() *CreateTitleStruct {
	request := CreateTitleStruct{Message: ""}
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
	titleStruct.Title = database.Title{
		Name:      generics.GetText(c, "name"),
		MainGroup: generics.GetText(c, "mainGroup"),
		SubGroup:  generics.GetText(c, "subGroup"),
		Flair:     gen.GetNullString(c, "flair"),
		Info: database.TitleInfo{
			Names: generics.GetStringArray(c, "user"),
		},
	}
	titleRef := &titleStruct.Title
	infoRef := &titleRef.Info

	switch true {
	case generics.CheckOrgOrTitle(titleStruct, titleRef):
	case gen.CheckAccountList(titleStruct, &infoRef.Names):
	case titleStruct.tryCreation():
	default:
		titleStruct.updateGroups()
		titleStruct.updateFlairs()
		titleStruct.refreshHierarchy()
		titleStruct.Message = generics.SuccessFullCreationTitle + "\n" + titleStruct.Message
	}

	return
}

func (titleStruct *CreateTitleStruct) tryCreation() bool {
	err := titleStruct.Title.CreateMe()
	if err != nil {
		titleStruct.Message = generics.TitleCreationError + "\n" + titleStruct.Message
		return true
	}
	return false
}

func (titleStruct *CreateTitleStruct) updateGroups() {
	if helper.GetPositionOfString(titleStruct.ExistingSubGroup, titleStruct.Title.SubGroup) == -1 {
		titleStruct.ExistingSubGroup = append(titleStruct.ExistingSubGroup, titleStruct.Title.SubGroup)
	}
	if helper.GetPositionOfString(titleStruct.ExistingMainGroup, titleStruct.Title.MainGroup) == -1 {
		titleStruct.ExistingMainGroup = append(titleStruct.ExistingMainGroup, titleStruct.Title.MainGroup)
	}
}

func (titleStruct *CreateTitleStruct) updateFlairs() {
	titleRef := &titleStruct.Title
	//TODO rework flair system
	if titleRef.Flair.Valid {
		err := dataLogic.UpdateFlairs([]string{}, titleRef.Info.Names, titleRef.Flair.String)
		if err != nil {
			titleStruct.Message = generics.FlairUpdateError + "\n" + titleStruct.Message
		}
	}
}

func (titleStruct *CreateTitleStruct) refreshHierarchy() {
	err := dataLogic.RefreshTitleHierarchy()
	if err != nil {
		titleStruct.Message = generics.RefresingTitleHierachyDidNotWork + "\n" + titleStruct.Message
	}
}
