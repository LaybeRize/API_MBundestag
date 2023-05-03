package htmlZwitscher

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ZwitscherListViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Zwitscherübersicht",
		Site:     "zwitscherList",
		Template: "zwitscherList",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ZwitscherSingleViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Zwitscher",
		Site:     "viewZwitscher",
		Template: "viewZwitscher",
	}
}
