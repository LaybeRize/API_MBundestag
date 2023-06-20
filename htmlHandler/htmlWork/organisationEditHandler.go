package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"errors"
	"github.com/gin-gonic/gin"
)

type EditOrganisationStruct struct {
	Organisation      dataLogic.Organisation
	OrgNames          []string
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	generics.MessageStruct
}

func getEmptyEditOrgStruct() *EditOrganisationStruct {
	request := EditOrganisationStruct{}
	htmlHandler.FillAllNotSuspendedNames(&request)
	htmlHandler.FillOrganisationGroups(&request)
	htmlHandler.FillOrganisationNames(&request)
	return &request
}

func GetEditOrganisationPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	err := errors.New("placeholder")
	editOrg := getEmptyEditOrgStruct()

	if !generics.GetIfEmptyQuery(c, "org") {
		err = editOrg.Organisation.GetMe(c.Query("org"))
	}

	if err != nil {
		editOrg.Organisation.Status = database.Public
		editOrg.Organisation.Member = []string{}
		editOrg.Organisation.Admins = []string{}
	}

	htmlHandler.MakeSite(editOrg, c, &acc)
}

func PostEditOrganisationPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	if generics.GetIfType(c, "search") {
		PostSearchOrganisationPage(c, &acc)
		return
	}

	htmlHandler.MakeSite(validateOrganisationEdit(c), c, &acc)
}

func PostSearchOrganisationPage(c *gin.Context, acc *database.Account) {
	htmlHandler.MakeSite(vaildateOrganisationSearch(c), c, acc)
}

func vaildateOrganisationSearch(c *gin.Context) (editOrg *EditOrganisationStruct) {

	editOrg = getEmptyEditOrgStruct()
	err := editOrg.Organisation.GetMe(generics.GetText(c, "name"))

	if err == nil {
		editOrg.Message = generics.SuccessfulFoundOrg + "\n" + editOrg.Message
		editOrg.Positiv = true
		return
	}
	editOrg.Message = generics.OrgFindingError + "\n" + editOrg.Message
	editOrg.Organisation = dataLogic.Organisation{
		Name:   generics.GetText(c, "name"),
		Status: database.Public,
		Member: []string{},
		Admins: []string{},
	}
	return
}

func validateOrganisationEdit(c *gin.Context) (orgStruct *EditOrganisationStruct) {
	orgStruct = getEmptyEditOrgStruct()
	orgStruct.Organisation = dataLogic.Organisation{
		Name:      generics.GetText(c, "name"),
		MainGroup: generics.GetText(c, "mainGroup"),
		SubGroup:  generics.GetText(c, "subGroup"),
		Flair:     generics.GetText(c, "flair"),
		Status:    database.StatusString(generics.GetText(c, "status")),
		Admins:    generics.GetStringArray(c, "admins"),
		Member:    generics.GetStringArray(c, "user"),
	}
	orgRef := &orgStruct.Organisation
	for _, str := range orgRef.Member {
		orgRef.Admins = help.RemoveFirstStringOccurrenceFromArray(orgRef.Admins, str)
	}

	switch true {
	case orgStruct.Message.CheckOrgOrTitle(orgRef):
	case orgStruct.checkIfExists():
	case orgStruct.Message.CheckOrgStatus(orgRef.Status):
	default:
		orgStruct.Organisation.ChangeMe(&orgStruct.Message, &orgStruct.Positiv)
	}
	return
}

func (orgStruct *EditOrganisationStruct) checkIfExists() bool {
	original := database.Organisation{}
	err := original.GetByName(orgStruct.Organisation.Name)
	if err != nil {
		orgStruct.Message = generics.OrgEditNonExistantElement + "\n" + orgStruct.Message
		return true
	}
	return false
}
