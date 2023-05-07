package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"strings"
)

type OrgansationNameEdit struct {
	OrganisationName string
	Names            []string
	User             []string
	Message          string
}

func GetOrganisationUserHandler(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	var org *database.Organisation
	b, org = checkIfUserHasAdminAccount(acc, strings.TrimSpace(c.Query("org")))
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	orgEdit := &OrgansationNameEdit{
		OrganisationName: org.Name,
		//User:             org.Info.User,
	}
	htmlHandler.FillAllNotSuspendedNames(orgEdit)
	htmlHandler.MakeSite(orgEdit, c, &acc)
}

var UserSuccessfullyChanged = "Nutzer erfolgreich verändert"
var ChangeUserOnOrgError = "Es ist ein Fehler beim verändern der Nutzer aufgetreten"
var FlairUserChangeError = "Es ist ein Fehler beim verändern der Flairs der Nutzer augetreten"

func PostOrganisationUserHandler(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	var org *database.Organisation
	b, org = checkIfUserHasAdminAccount(acc, generics.GetText(c, "name"))
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	orgEdit := &OrgansationNameEdit{
		OrganisationName: org.Name,
		User:             generics.GetStringArray(c, "user"),
	}
	htmlHandler.FillAllNotSuspendedNames(orgEdit)

	//old := make([]string, len(org.Info.User))
	//copy(old, org.Info.User)

	switch true {
	//case gen.CheckAccountList(orgEdit, &orgEdit.User):
	case orgEdit.updateAllowed(org):
	case orgEdit.trySaving(org):
	default:
		//orgEdit.tryChangingFlair(org, old)
		orgEdit.Message = UserSuccessfullyChanged + "\n" + orgEdit.Message
	}

	htmlHandler.MakeSite(orgEdit, c, &acc)
}

func checkIfUserHasAdminAccount(acc database.Account, name string) (b bool, org *database.Organisation) {
	/*b = false
	org = &database.Organisation{}
	err := org.GetByName(name)
	if err != nil {
		return
	}
	if acc.Role == database.HeadAdmin || acc.Role == database.Admin {
		return true, org
	}
	list := database.AccountList{}
	err = list.GetAllPressAccountsFromAccountPlusSelf(acc)
	if err != nil {
		return
	}
	for _, p := range list {
		if helper.GetPositionOfString(org.Info.Admins, p.DisplayName) != -1 {
			return true, org
		}
	}*/
	return
}

func (e *OrgansationNameEdit) updateAllowed(org *database.Organisation) bool {
	/*for _, name := range org.Info.Admins {
		e.User = helper.RemoveFirstStringOccurrenceFromArray(e.User, name)
	}
	org.Info.User = e.User

	acc := database.Account{}
	org.Info.Viewer = []string{}
	for _, str := range org.Info.User {
		err := acc.GetByDisplayName(str)
		if acc.Role == database.PressAccount {
			err = acc.GetByID(acc.Linked.Int64)
		}
		if err != nil {
			e.Message = generics.ViewerError + "\n" + e.Message
			return true
		}
		org.Info.Viewer = append(org.Info.Viewer, acc.DisplayName)
	}
	for _, str := range org.Info.Admins {
		err := acc.GetByDisplayName(str)
		if acc.Role == database.PressAccount {
			err = acc.GetByID(acc.Linked.Int64)
		}
		if err != nil {
			e.Message = generics.ViewerError + "\n" + e.Message
			return true
		}
		org.Info.Viewer = append(org.Info.Viewer, acc.DisplayName)
	}
	org.Info.Viewer = helper.RemoveDuplicates(org.Info.Viewer)*/
	return false
}

func (e *OrgansationNameEdit) trySaving(org *database.Organisation) bool {
	err := org.SaveChanges()
	if err != nil {
		e.Message = ChangeUserOnOrgError + "\n" + e.Message
		return true
	}
	return false
}

func (e *OrgansationNameEdit) tryChangingFlair(org *database.Organisation, old []string) {
	/*if org.Flair.Valid {
		err := dataLogic.UpdateFlairs(old, org.Info.User, org.Flair.String)
		if err != nil {
			e.Message = FlairUserChangeError + "\n" + e.Message
		}
	}*/
}
