package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/crypto/openpgp"
)

type CreateOrganisationStruct struct {
	Organisation      dataLogic.Organisation
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	generics.MessageStruct
}

func getEmptyCreateOrgStruct() *CreateOrganisationStruct {
	request := CreateOrganisationStruct{}
	htmlHandler.FillAllNotSuspendedNames(&request)
	htmlHandler.FillOrganisationGroups(&request)
	return &request
}

func GetCreateOrganisationPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	orgStruct := getEmptyCreateOrgStruct()
	orgStruct.Organisation = dataLogic.Organisation{
		Status: database.Public,
		Admins: []string{},
		Member: []string{},
	}
	htmlHandler.MakeSite(orgStruct, c, &acc)
}

func PostCreateOrganisationPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(validateOrganisationCreate(c), c, &acc)
}

func validateOrganisationCreate(c *gin.Context) (orgStruct *CreateOrganisationStruct) {
	orgStruct = getEmptyCreateOrgStruct()
	orgStruct.Organisation = dataLogic.Organisation{
		Name:      generics.GetText(c, "name"),
		MainGroup: generics.GetText(c, "mainGroup"),
		SubGroup:  generics.GetText(c, "subGroup"),
		Flair:     generics.GetText(c, "flair"),
		Status:    database.StatusString(generics.GetText(c, "status")),
		Admins:    generics.GetStringArray(c, "admins"),
		Member:    generics.GetStringArray(c, "user"),
	}
	//easier access to the org info
	orgRef := &orgStruct.Organisation
	//remove any user also added as admins
	for _, str := range orgRef.Member {
		orgStruct.Organisation.Admins = help.RemoveFirstStringOccurrenceFromArray(orgStruct.Organisation.Admins, str)
	}

	switch true {
	case orgStruct.Message.CheckOrgOrTitle(orgRef):
	case orgStruct.Message.CheckOrgStatus(orgRef.Status):
	default:
		orgStruct.Organisation.CreateMe(&orgStruct.Message, &orgStruct.Positiv)
	}
	return
}
