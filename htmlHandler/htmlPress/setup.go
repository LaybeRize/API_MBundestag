package htmlPress

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(CreateArticleStruct{})] = htmlHandler.BasicStruct{
		Title:    "Artikel einreichen",
		Site:     "createArticle",
		Template: "createArticle",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(NewspaperListViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Zeitungsübersicht",
		Site:     "newspaperList",
		Template: "newspaperList",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(NewspaperHiddenListViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Zeitungsübersicht",
		Site:     "newspaperAdminList",
		Template: "newspaperList",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(PublicationViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Zeitungsansicht",
		Site:     "viewPublication",
		Template: "viewPublication",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(RejectArticleStruct{})] = htmlHandler.BasicStruct{
		Title:    "Artikel ablehnen",
		Site:     "rejectArticle",
		Template: "rejectArticle",
	}
}
