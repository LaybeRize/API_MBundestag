package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	gen "API_MBundestag/htmlHandler/generics"
	"API_MBundestag/htmlHandler/htmlBasics"
	"errors"
	"github.com/gin-gonic/gin"
)

type EditTitleStruct struct {
	TitleNames        []string
	Title             database.Title
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	Message           string
}

func getEmptyEditTitleStruct() *EditTitleStruct {
	request := EditTitleStruct{Message: ""}
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
	err := errors.New("placeholder")

	titleStruct := getEmptyEditTitleStruct()
	if !htmlHandler.GetIfEmptyQuery(c, "title") {
		err = titleStruct.Title.GetByName(c.Query("title"))
	}

	if err != nil {
		titleStruct.Title = database.Title{Info: database.TitleInfo{Names: []string{}}}
	}

	htmlHandler.MakeSite(titleStruct, c, &acc)
}

func PostEditTitlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	if htmlHandler.GetIfType(c, "search") {
		htmlHandler.MakeSite(validateSearchTitle(c), c, &acc)
		return
	}
	if htmlHandler.GetIfType(c, "delete") {
		htmlHandler.MakeSite(validateDeleteTitle(c), c, &acc)
		return
	}

	htmlHandler.MakeSite(validateEditTitle(c), c, &acc)
}

func validateSearchTitle(c *gin.Context) (editStruct *EditTitleStruct) {
	editStruct = getEmptyEditTitleStruct()
	err := editStruct.Title.GetByName(htmlHandler.GetText(c, "name"))
	if err != nil {
		editStruct.Message = generics.TitleDoesNotExists + "\n" + editStruct.Message
	} else {
		editStruct.Message = generics.SuccessFullFoundTitle + "\n" + editStruct.Message
	}
	return
}

func validateDeleteTitle(c *gin.Context) (editStruct *EditTitleStruct) {
	editStruct = getEmptyEditTitleStruct()
	err := editStruct.Title.GetByName(htmlHandler.GetText(c, "name"))
	if err != nil {
		editStruct.Message = generics.TitleDoesNotExists + "\n" + editStruct.Message
		return
	}
	err = editStruct.Title.DeleteMe()
	if err != nil {
		editStruct.Message = generics.ErrorWhileDeletingTitle + "\n" + editStruct.Message
		return
	}
	editStruct.TitleNames = helper.RemoveFirstStringOccurrenceFromArray(editStruct.TitleNames, editStruct.Title.Name)

	err = dataLogic.RefreshTitleHierarchy()
	if err != nil {
		editStruct.Message = generics.RefresingTitleHierachyDidNotWork + "\n" + editStruct.Message
	}
	editStruct.Message = generics.SuccesfulDeletedTitle + "\n" + editStruct.Message

	return
}

func validateEditTitle(c *gin.Context) (editStruct *EditTitleStruct) {
	editStruct = getEmptyEditTitleStruct()
	editStruct.Title = database.Title{
		Name:      htmlHandler.GetText(c, "newName"),
		MainGroup: htmlHandler.GetText(c, "mainGroup"),
		SubGroup:  htmlHandler.GetText(c, "subGroup"),
		Flair:     gen.GetNullString(c, "flair"),
		Info: database.TitleInfo{
			Names: htmlHandler.GetStringArray(c, "user"),
		},
	}
	titleRef := &editStruct.Title
	infoRef := &titleRef.Info
	oldTitle := &database.Title{}

	switch true {
	case editStruct.getOldTitle(c, oldTitle):
	case htmlHandler.CheckOrgOrTitle(editStruct, titleRef):
		editStruct.Title.Name = oldTitle.Name
	case gen.CheckAccountList(editStruct, &infoRef.Names):
		editStruct.Title.Name = oldTitle.Name
	case editStruct.tryChange(oldTitle.Name):
	default:
		editStruct.updateGroups()
		editStruct.changeFlair(oldTitle)
		editStruct.refreshHierarchy()
		editStruct.Message = generics.SuccessFullEditTitle + "\n" + editStruct.Message
	}
	return
}

func (editStruct *EditTitleStruct) getOldTitle(c *gin.Context, oldTitle *database.Title) bool {
	err := oldTitle.GetByName(htmlHandler.GetText(c, "name"))
	if err != nil {
		editStruct.Title.Name = htmlHandler.GetText(c, "name")
		editStruct.Message = generics.TitleDoesNotExists + "\n" + editStruct.Message
		return true
	}
	return false
}

func (editStruct *EditTitleStruct) tryChange(oldName string) bool {
	err := editStruct.Title.ChangeTitleName(oldName)
	if err != nil {
		editStruct.Message = generics.TitelUpdateError + "\n" + editStruct.Message
		editStruct.Title.Name = oldName
		return true
	}
	return false
}

func (editStruct *EditTitleStruct) updateGroups() {
	if helper.GetPositionOfString(editStruct.ExistingSubGroup, editStruct.Title.SubGroup) == -1 {
		editStruct.ExistingSubGroup = append(editStruct.ExistingSubGroup, editStruct.Title.SubGroup)
	}
	if helper.GetPositionOfString(editStruct.ExistingMainGroup, editStruct.Title.MainGroup) == -1 {
		editStruct.ExistingMainGroup = append(editStruct.ExistingMainGroup, editStruct.Title.MainGroup)
	}
}

func (editStruct *EditTitleStruct) changeFlair(old *database.Title) {
	newTitle := &editStruct.Title
	var err error
	switch true {
	case !old.Flair.Valid && !newTitle.Flair.Valid:
	case !old.Flair.Valid && newTitle.Flair.Valid:
		err = dataLogic.UpdateFlairs([]string{}, newTitle.Info.Names, newTitle.Flair.String)
	case old.Flair.Valid && !newTitle.Flair.Valid:
		err = dataLogic.UpdateFlairs(old.Info.Names, []string{}, old.Flair.String)
	case old.Flair.String == newTitle.Flair.String:
		err = dataLogic.UpdateFlairs(old.Info.Names, newTitle.Info.Names, old.Flair.String)
	default:
		err = dataLogic.UpdateFlairs(old.Info.Names, newTitle.Info.Names, old.Flair.String, newTitle.Flair.String)
	}
	if err != nil {
		editStruct.Message = generics.FlairUpdateError + "\n" + editStruct.Message
	}
}

func (editStruct *EditTitleStruct) refreshHierarchy() {
	err := dataLogic.RefreshTitleHierarchy()
	if err != nil {
		editStruct.Message = generics.RefresingTitleHierachyDidNotWork + "\n" + editStruct.Message
	}
}
