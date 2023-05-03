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

type EditOrganisationStruct struct {
	Organisation      database.Organisation
	OrgNames          []string
	ExistingMainGroup []string
	ExistingSubGroup  []string
	Names             []string
	Message           string
}

func getEmptyEditOrgStruct() *EditOrganisationStruct {
	request := EditOrganisationStruct{Message: ""}
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
	org := database.Organisation{}

	if !generics.GetIfEmptyQuery(c, "org") {
		err = org.GetByName(c.Query("org"))
	}

	if err != nil {
		org.Status = database.Public
		org.Info.Admins = []string{}
		org.Info.User = []string{}
	}

	editOrg := getEmptyEditOrgStruct()
	editOrg.Organisation = org
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
	org := database.Organisation{}
	err := org.GetByName(generics.GetText(c, "name"))

	editOrg = getEmptyEditOrgStruct()

	if err == nil {
		editOrg.Message = generics.SuccessFullFindOrg + "\n" + editOrg.Message
		editOrg.Organisation = org
		return
	}
	editOrg.Message = generics.OrgFindingError + "\n" + editOrg.Message
	editOrg.Organisation = database.Organisation{
		Name:   generics.GetText(c, "name"),
		Status: database.Public,
		Info: database.OrganisationInfo{
			Admins: []string{},
			User:   []string{},
		},
	}
	return
}

func validateOrganisationEdit(c *gin.Context) (orgStruct *EditOrganisationStruct) {
	orgStruct = getEmptyEditOrgStruct()
	orgStruct.Organisation = database.Organisation{
		Name:      generics.GetText(c, "name"),
		MainGroup: generics.GetText(c, "mainGroup"),
		SubGroup:  generics.GetText(c, "subGroup"),
		Flair:     gen.GetNullString(c, "flair"),
		Status:    database.StatusString(generics.GetText(c, "status")),
		Info: database.OrganisationInfo{
			Admins: generics.GetStringArray(c, "admins"),
			User:   generics.GetStringArray(c, "user"),
			Viewer: []string{},
		},
	}
	orgRef := &orgStruct.Organisation
	infoRef := &orgRef.Info
	for _, str := range infoRef.User {
		infoRef.Admins = helper.RemoveFirstStringOccurrenceFromArray(infoRef.Admins, str)
	}

	original := &database.Organisation{}
	switch true {
	case generics.CheckOrgOrTitle(orgStruct, orgRef):
	case orgStruct.checkIfExists(original):
	case gen.CheckAccountList(orgStruct, &infoRef.Admins):
	case gen.CheckAccountList(orgStruct, &infoRef.User):
	case generics.CheckOrgStatus(orgStruct, orgRef.Status):
	case orgStruct.addViewer(infoRef.Admins):
	case orgStruct.addViewer(infoRef.User):
	case orgStruct.tryCreation():
	default:
		orgStruct.updateGroups()
		orgStruct.tryUpdatingFlairs(original)
		orgStruct.Message = generics.SuccessFullChangeOrg + "\n" + orgStruct.Message
	}
	return
}

func (orgStruct *EditOrganisationStruct) checkIfExists(original *database.Organisation) bool {
	err := original.GetByName(orgStruct.Organisation.Name)
	if err != nil {
		orgStruct.Message = generics.OrgEditNonExistantElement + "\n" + orgStruct.Message
		return true
	}
	return false
}

func (orgStruct *EditOrganisationStruct) addViewer(array []string) bool {
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

func (orgStruct *EditOrganisationStruct) tryCreation() bool {
	orgRef := &orgStruct.Organisation
	infoRef := &orgRef.Info
	//invalidate string if the status is secret
	if orgRef.Status == database.Secret {
		orgRef.Flair.Valid = false
	}

	//empty organisation, if it's hidden and make the flair unavailable
	if orgRef.Status == database.Hidden {
		orgRef.Flair.String = ""
		orgRef.Flair.Valid = false
		infoRef.User = []string{}
		infoRef.Admins = []string{}
		infoRef.Viewer = []string{}
	}

	//make sure that the org is correctly saved
	err := orgRef.SaveChanges()
	if err != nil {
		orgStruct.Message = generics.OrganisationEditError + "\n" + orgStruct.Message
		return true
	}
	return false
}

func (orgStruct *EditOrganisationStruct) updateGroups() {
	if helper.GetPositionOfString(orgStruct.ExistingSubGroup, orgStruct.Organisation.SubGroup) == -1 {
		orgStruct.ExistingSubGroup = append(orgStruct.ExistingSubGroup, orgStruct.Organisation.SubGroup)
	}
	if helper.GetPositionOfString(orgStruct.ExistingMainGroup, orgStruct.Organisation.MainGroup) == -1 {
		orgStruct.ExistingMainGroup = append(orgStruct.ExistingMainGroup, orgStruct.Organisation.MainGroup)
	}
}

func (orgStruct *EditOrganisationStruct) tryUpdatingFlairs(original *database.Organisation) {
	org := &orgStruct.Organisation
	var err error
	var err2 error
	//TODO rework flair system
	switch true {
	case !org.Flair.Valid && !original.Flair.Valid:
	case org.Flair.Valid && !original.Flair.Valid:
		err = dataLogic.UpdateFlairs([]string{}, org.Info.User, org.Flair.String)
		err2 = dataLogic.UpdateFlairs([]string{}, org.Info.Admins, org.Flair.String)
	case !org.Flair.Valid && original.Flair.Valid:
		err = dataLogic.UpdateFlairs(original.Info.User, []string{}, original.Flair.String)
		err2 = dataLogic.UpdateFlairs(original.Info.Admins, []string{}, original.Flair.String)
	case org.Flair.String == original.Flair.String:
		err = dataLogic.UpdateFlairs(original.Info.User, org.Info.User, org.Flair.String)
		err2 = dataLogic.UpdateFlairs(original.Info.Admins, org.Info.Admins, org.Flair.String)
	default:
		err = dataLogic.UpdateFlairs(original.Info.User, org.Info.User, original.Flair.String, org.Flair.String)
		err2 = dataLogic.UpdateFlairs(original.Info.Admins, org.Info.Admins, original.Flair.String, org.Flair.String)
	}

	if err != nil || err2 != nil {
		orgStruct.Message = generics.FlairUpdateError + "\n" + orgStruct.Message
	}
}
