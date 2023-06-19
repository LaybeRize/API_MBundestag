package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type OrgansationNameEdit struct {
	OrganisationName string
	Names            []string
	User             []string
	generics.MessageStruct
}

func GetOrganisationUserHandler(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	org := dataLogic.Organisation{}
	err := org.GetMeWhenAdmin(c.Query("name"), acc.ID)
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.CouldNotFindOrganisation)
		return
	}

	orgEdit := &OrgansationNameEdit{
		OrganisationName: org.Name,
		User:             org.Member,
		MessageStruct: generics.MessageStruct{
			Message: generics.SuccessFullFindOrg,
			Positiv: true,
		},
	}
	htmlHandler.FillAllNotSuspendedNames(orgEdit)

	htmlHandler.MakeSite(orgEdit, c, &acc)
}

func PostOrganisationUserHandler(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	org := dataLogic.Organisation{}
	err := org.GetMeWhenAdmin(generics.GetText(c, "name"), acc.ID)
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.CouldNotFindOrganisation)
		return
	}

	org.Member = generics.GetStringArray(c, "user")
	for _, str := range org.Admins {
		org.Member = help.RemoveFirstStringOccurrenceFromArray(org.Member, str)
	}

	orgEdit := &OrgansationNameEdit{
		OrganisationName: org.Name,
		User:             org.Member,
	}
	htmlHandler.FillAllNotSuspendedNames(orgEdit)

	org.ChangeOnlyMembers(&orgEdit.Message, &orgEdit.Positiv)

	htmlHandler.MakeSite(orgEdit, c, &acc)
}
