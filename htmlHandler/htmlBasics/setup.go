package htmlBasics

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ErrorCode{})] = htmlHandler.BasicStruct{
		Title:    "Fehlerseite",
		Site:     "",
		Template: "error",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(StartPageStruct{})] = htmlHandler.BasicStruct{
		Title:    "Startseite",
		Site:     "start",
		Template: "start",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(TestPreviewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Vorschautest",
		Site:     "testPreview",
		Template: "testPreview",
	}

	htmlHandler.AddFunctionToLinks("/start", GetStartPage)
	htmlHandler.AddFunctionToLinks("/start", PostStartPage)
	htmlHandler.AddFunctionToLinks("/markdown", PostJsonMarkdown)
}
