package htmlAccount

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ListUserStruct{})] = htmlHandler.BasicStruct{
		Title:    "Nutzerliste",
		Site:     "adminViewUser",
		Template: "adminViewUser",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(EditUserStruct{})] = htmlHandler.BasicStruct{
		Title:    "Nutzerbearbeitung",
		Site:     "editUser",
		Template: "editUser",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(CreateUserStruct{})] = htmlHandler.BasicStruct{
		Title:    "Nutzererstellung",
		Site:     "createUser",
		Template: "createUser",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ViewUserInfoStruct{})] = htmlHandler.BasicStruct{
		Title:    "Meine Übersicht",
		Site:     "viewPersonalInfo",
		Template: "viewPersonalInfo",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(PasswordChangeStruct{})] = htmlHandler.BasicStruct{
		Title:    "Passwort ändern",
		Site:     "password",
		Template: "password",
	}
}
