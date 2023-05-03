package htmlLetter

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ViewSingleLetter{})] = htmlHandler.BasicStruct{
		Title:    "Briefansicht",
		Site:     "viewLetter",
		Template: "viewLetter",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(LetterCreatePageStruct{})] = htmlHandler.BasicStruct{
		Title:    "Brief erstellen",
		Site:     "createLetter",
		Template: "createLetter",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ModMailCreatePageStruct{})] = htmlHandler.BasicStruct{
		Title:    "Moderationsbrief erstellen",
		Site:     "createModMail",
		Template: "createLetter",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(AdminLetterViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Brief suchen",
		Site:     "adminViewLetter",
		Template: "adminViewLetter",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ViewModMailListStrcut{})] = htmlHandler.BasicStruct{
		Title:    "Moderationsbriefe",
		Site:     "modMailList",
		Template: "letterList",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ViewLetterListStruct{})] = htmlHandler.BasicStruct{
		Title:    "Deine Briefe",
		Site:     "letterList",
		Template: "letterList",
	}
}
