package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type HiddenOrganisationStruct dataLogic.OrganisationMainGroupArray

func GetOrganisationViewPage(c *gin.Context) {
	acc, _ := dataLogic.CheckUserPrivileged(c)

	orgHierarchy := dataLogic.OrganisationMainGroupArray{}
	err := orgHierarchy.GetOrganisationHierarchy(acc, false)

	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLodingOrganisationView)
		return
	}

	htmlHandler.MakeSite(&orgHierarchy, c, &acc)
}

func GetHiddenOrganisationViewPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	orgHierarchy := dataLogic.OrganisationMainGroupArray{}
	err := orgHierarchy.GetOrganisationHierarchy(acc, true)

	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLodingOrganisationView)
		return
	}
	orgHierarchy.SetAmountForMainGroup()

	h := HiddenOrganisationStruct(orgHierarchy)
	htmlHandler.MakeSite(&h, c, &acc)
}
