package htmlDocuments

import (
	"API_MBundestag/htmlHandler"
)

func Setup() {
	htmlHandler.PageIdentityMap[htmlHandler.Identity(PostViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Legislativer Text",
		Site:     "viewPost",
		Template: "viewPost",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(DiscussionViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Diskussion",
		Site:     "viewDiscussion",
		Template: "viewDiscussion",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(VoteViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Abstimmung",
		Site:     "viewVote",
		Template: "viewVote",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(DiscussionCreateStruct{})] = htmlHandler.BasicStruct{
		Title:    "Diskussion erstellen",
		Site:     "createDiscussion",
		Template: "createDiscussion",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(CreatePostPageStruct{})] = htmlHandler.BasicStruct{
		Title:    "Legislativen Text erstellen",
		Site:     "createPost",
		Template: "createPost",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(VoteCreateStruct{})] = htmlHandler.BasicStruct{
		Title:    "Abstimmung erstellen",
		Site:     "createVote",
		Template: "createVote",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(DocumentNavigationStruct{})] = htmlHandler.BasicStruct{
		Title:    "Dokumenterstellernavigation",
		Site:     "postNavigation",
		Template: "postNavigation",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ViewDocumentListStruct{})] = htmlHandler.BasicStruct{
		Title:    "Dokumentübersicht",
		Site:     "listDocuments",
		Template: "listDocuments",
	}

	htmlHandler.PageIdentityMap[htmlHandler.Identity(SingleVoteViewStruct{})] = htmlHandler.BasicStruct{
		Title:    "Abstimmung",
		Site:     "singleVote",
		Template: "singleVote",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(MultipleOptionsVoteStruct{})] = htmlHandler.BasicStruct{
		Title:    "Abstimmung",
		Site:     "multipleVote",
		Template: "multipleVote",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(RankingVoteStruct{})] = htmlHandler.BasicStruct{
		Title:    "Abstimmung",
		Site:     "rankedVote",
		Template: "rankedVote",
	}
	htmlHandler.PageIdentityMap[htmlHandler.Identity(ForAgainstVoteStruct{})] = htmlHandler.BasicStruct{
		Title:    "Abstimmung",
		Site:     "threeChoicesVote",
		Template: "threeChoicesVote",
	}
}
