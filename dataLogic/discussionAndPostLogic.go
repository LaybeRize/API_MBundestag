package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help"
	"time"
)

var DiscussioNotFoundError help.Message = "Diskussion konnte nicht gefunden werden"
var DiscussionAlreadyClosed help.Message = "Diese Diskussion ist bereits beendet"
var CanNotCommentOnClosedDiscussion help.Message = "Du kannst keine Kommentare unter einer geschlossenen Diskussion schreiben"
var CouldNotSaveComment help.Message = "Kommentar konnte nicht gespeichert werden"
var SuccessfulComment help.Message = "Kommentar wurde gespeichert"

func AddComment(discussion *database.Document, comment *database.Discussions, msg *help.Message, positiv *bool) {
	documentLock.Lock()
	defer documentLock.Unlock()

	err := discussion.GetByID(discussion.UUID)
	if err != nil {
		*msg = DiscussioNotFoundError + "\n" + *msg
		return
	}

	if discussion.Type == database.RunningDiscussion && discussion.Info.Finishing.UTC().Before(time.Now().UTC()) {
		go CloseDiscussionOrVote(discussion.UUID)
		*msg = DiscussionAlreadyClosed + "\n" + *msg
		return
	} else if discussion.Type != database.RunningDiscussion {
		*msg = CanNotCommentOnClosedDiscussion + "\n" + *msg
		return
	}

	discussion.Info.Discussion = append(discussion.Info.Discussion, *comment)
	err = discussion.SaveChanges()
	if err != nil {
		*msg = CouldNotSaveComment + "\n" + *msg
		return
	}
	*positiv = true
	*msg = SuccessfulComment + "\n" + *msg
}

var CommentNotFoundError help.Message = "Kommentar konnte nicht gefunden werden"
var CommentSuccessfulHidden help.Message = "Kommentar wurde versteckt"
var CommentSuccesfulUnhidden help.Message = "Kommentar wurde erfolgreich freigeschaltet"
var CouldNotSaveCommentHidden help.Message = "Kommentarveränderung konnte nicht gespeichert werden"

func ToggleHideDiscussionComment(discussion *database.Document, commentID string, msg *help.Message, positiv *bool) {
	documentLock.Lock()
	defer documentLock.Unlock()

	err := discussion.GetByID(discussion.UUID)
	if err != nil {
		*msg = DiscussioNotFoundError + "\n" + *msg
		return
	}
	pos := -1
	for i, comment := range discussion.Info.Discussion {
		if comment.UUID == commentID {
			pos = i
		}
	}
	if pos == -1 {
		*msg = CommentNotFoundError + "\n" + *msg
		return
	}
	discussion.Info.Discussion[pos].Hidden = !discussion.Info.Discussion[pos].Hidden
	err = discussion.SaveChanges()
	if err != nil {
		*msg = CouldNotSaveCommentHidden + "\n" + *msg
		return
	}
	*positiv = true
	if discussion.Info.Discussion[pos].Hidden {
		*msg = CommentSuccessfulHidden + "\n" + *msg
	} else {
		*msg = CommentSuccesfulUnhidden + "\n" + *msg
	}
}

var DocumentNotFoundError help.Message = "Dokument konnte nicht gefunden werden"
var DocumentIsNotLegislativeText help.Message = "Das Dokument besitzt keine Tags"
var CouldNotSavePostTag help.Message = "Tag konnte nicht gespeichert werden"
var SuccessfulPostTag help.Message = "Tag wurde erfolgreichspeichert"

func AddPostTag(postTag *database.Document, post *database.Posts, msg *help.Message, positiv *bool) {
	documentLock.Lock()
	defer documentLock.Unlock()

	err := postTag.GetByID(postTag.UUID)
	if err != nil {
		*msg = DocumentNotFoundError + "\n" + *msg
		return
	}
	if postTag.Type != database.LegislativeText {
		*msg = DocumentIsNotLegislativeText + "\n" + *msg
		return
	}
	postTag.Info.Post = append(postTag.Info.Post, *post)
	err = postTag.SaveChanges()
	if err != nil {
		*msg = CouldNotSavePostTag + "\n" + *msg
		return
	}
	*positiv = true
	*msg = SuccessfulPostTag + "\n" + *msg
}

var CouldNotFindPostTag help.Message = "Tag konnte nicht gefunden werden"
var CouldNotChangePostTag help.Message = "Tag konnte nicht geändert werden"
var SuccessfulPostTagHidden help.Message = "Tag wurde erfolgreich versteckt"
var SuccessfulPostTagUnhidden help.Message = "Tag wurde erfolgreich freigeschaltet"

func ToggleHidePostTag(postTag *database.Document, postTagID string, msg *help.Message, positiv *bool) {
	documentLock.Lock()
	defer documentLock.Unlock()

	err := postTag.GetByID(postTag.UUID)
	if err != nil {
		*msg = DocumentNotFoundError + "\n" + *msg
		return
	}
	pos := -1
	for i, post := range postTag.Info.Post {
		if post.UUID == postTagID {
			pos = i
		}
	}
	if pos == -1 {
		*msg = CouldNotFindPostTag + "\n" + *msg
		return
	}
	postTag.Info.Post[pos].Hidden = !postTag.Info.Post[pos].Hidden

	err = postTag.SaveChanges()
	if err != nil {
		*msg = CouldNotChangePostTag + "\n" + *msg
		return
	}
	*positiv = true
	if postTag.Info.Post[pos].Hidden {
		*msg = SuccessfulPostTagHidden + "\n" + *msg
	} else {
		*msg = SuccessfulPostTagUnhidden + "\n" + *msg
	}
}
