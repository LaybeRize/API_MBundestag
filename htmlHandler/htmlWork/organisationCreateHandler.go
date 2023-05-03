package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	gen "API_MBundestag/htmlHandler/generics"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/crypto/openpgp"
)

type CreateOrganisationStruct struct {
	Organisation      database.Organisation
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	Message           string
}

func getEmptyCreateOrgStruct() *CreateOrganisationStruct {
	request := CreateOrganisationStruct{Message: ""}
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
	orgStruct.Organisation = database.Organisation{
		Status: database.Public,
		Info: database.OrganisationInfo{
			Admins: []string{},
			User:   []string{},
		},
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
	orgStruct.Organisation = database.Organisation{
		Name:      htmlHandler.GetText(c, "name"),
		MainGroup: htmlHandler.GetText(c, "mainGroup"),
		SubGroup:  htmlHandler.GetText(c, "subGroup"),
		Flair:     gen.GetNullString(c, "flair"),
		Status:    database.StatusString(htmlHandler.GetText(c, "status")),
		Info: database.OrganisationInfo{
			Admins: htmlHandler.GetStringArray(c, "admins"),
			User:   htmlHandler.GetStringArray(c, "user"),
			Viewer: []string{},
		},
	}
	//easier access to the org info
	orgRef := &orgStruct.Organisation
	infoRef := &orgRef.Info
	//remove any user also added as admins
	for _, str := range infoRef.User {
		infoRef.Admins = helper.RemoveFirstStringOccurrenceFromArray(infoRef.Admins, str)
	}

	switch true {
	case htmlHandler.CheckOrgOrTitle(orgStruct, orgRef):
	case gen.CheckAccountList(orgStruct, &infoRef.Admins):
	case gen.CheckAccountList(orgStruct, &infoRef.User):
	case htmlHandler.CheckOrgStatus(orgStruct, orgRef.Status):
	case orgStruct.addViewer(infoRef.Admins):
	case orgStruct.addViewer(infoRef.User):
	case orgStruct.tryCreation():
	default:
		orgStruct.updateGroups()
		orgStruct.tryUpdatingFlairs()
		orgStruct.Message = generics.SuccessFullCreationOrg + "\n" + orgStruct.Message
	}
	return
}

func (orgStruct *CreateOrganisationStruct) addViewer(array []string) bool {
	infoRef := &orgStruct.Organisation.Info
	acc := database.Account{}
	for _, str := range array {
		err := acc.GetByDisplayName(str)
		if acc.Role == database.PressAccount {
			err = acc.GetByID(acc.Linked.Int64)
		}
		if err != nil {
			orgStruct.Message = generics.ViewerError + "\n" + orgStruct.Message
			return true
		}
		infoRef.Viewer = append(infoRef.Viewer, acc.DisplayName)
	}
	infoRef.Viewer = helper.RemoveDuplicates(infoRef.Viewer)
	return false
}

func (orgStruct *CreateOrganisationStruct) tryCreation() bool {
	orgRef := &orgStruct.Organisation
	//hidden and secret organisation are not allowed to have a flair
	if orgRef.Status == database.Secret || orgRef.Status == database.Hidden {
		orgRef.Flair.Valid = false
	}

	err := orgRef.CreateMe()
	//if anyhting goes wrong while creating the org
	if err != nil {
		orgStruct.Message = generics.OrganisationCreationError + "\n" + orgStruct.Message
		return true
	}
	return false
}

func (orgStruct *CreateOrganisationStruct) updateGroups() {
	//add to the existing main groups (because they exist now) if they are not already in the list
	if helper.GetPositionOfString(orgStruct.ExistingSubGroup, orgStruct.Organisation.SubGroup) == -1 {
		orgStruct.ExistingSubGroup = append(orgStruct.ExistingSubGroup, orgStruct.Organisation.SubGroup)
	}
	if helper.GetPositionOfString(orgStruct.ExistingMainGroup, orgStruct.Organisation.MainGroup) == -1 {
		orgStruct.ExistingMainGroup = append(orgStruct.ExistingMainGroup, orgStruct.Organisation.MainGroup)
	}
}

func (orgStruct *CreateOrganisationStruct) tryUpdatingFlairs() {
	orgRef := &orgStruct.Organisation
	infoRef := &orgRef.Info
	var err error
	var err2 error
	//TODO rework flair system
	if orgRef.Flair.Valid {
		err = dataLogic.UpdateFlairs([]string{}, infoRef.Admins, orgRef.Flair.String)
		err2 = dataLogic.UpdateFlairs([]string{}, infoRef.User, orgRef.Flair.String)
	}
	if err != nil || err2 != nil {
		orgStruct.Message = generics.FlairUpdateError + "\n" + orgStruct.Message
		return
	}
}
