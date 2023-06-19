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

	htmlHandler.AddFunctionToLinks("/create-user", GetCreateUserPage)
	htmlHandler.AddFunctionToLinks("/create-user", PostCreateUserPage)
	htmlHandler.AddFunctionToLinks("/edit-user", GetEditUserPage)
	htmlHandler.AddFunctionToLinks("/edit-user", PostEditUserPage)
	htmlHandler.AddFunctionToLinks("/view-user", GetAdminListUserPage)
	htmlHandler.AddFunctionToLinks("/self-info", GetViewOfProfilePage)
	htmlHandler.AddFunctionToLinks("/password", GetPasswordChangePage)
	htmlHandler.AddFunctionToLinks("/password", PostPasswordChangePage)
}
