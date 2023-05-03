package htmlWork

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(HiddenOrganisationStruct{})] = htmlHandler.BasicStruct{
		Title:    "Organisationsübersicht",
		Site:     "hiddenOrganisation",
		Template: "hiddenOrganisation",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(CreateOrganisationStruct{})] = htmlHandler.BasicStruct{
		Title:    "Organisationen erstellen",
		Site:     "createOrganisation",
		Template: "createOrganisation",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(EditOrganisationStruct{})] = htmlHandler.BasicStruct{
		Title:    "Organisationen bearbeiten",
		Site:     "editOrganisation",
		Template: "editOrganisation",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(CreateTitleStruct{})] = htmlHandler.BasicStruct{
		Title:    "Titel erstellen",
		Site:     "createTitle",
		Template: "createTitle",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(EditTitleStruct{})] = htmlHandler.BasicStruct{
		Title:    "Titel bearbeiten und löschen",
		Site:     "editTitle",
		Template: "editTitle",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(OrgansationNameEdit{})] = htmlHandler.BasicStruct{
		Title:    "Organisationsnutzer anpasssen",
		Site:     "editUserOrganisation",
		Template: "editUserOrganisation",
	}
}
