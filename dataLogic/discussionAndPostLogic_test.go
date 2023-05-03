package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDiscussionAndPostLogic(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupDiscussionAndPost", testSetupDiscussionAndPost)
	t.Run("testAddComment", testAddComment)
	t.Run("testHideComment", testHideComment)
	t.Run("testAddTag", testAddTag)
	t.Run("testHideTag", testHideTag)
}

func testHideTag(t *testing.T) {
	disc := database.Document{UUID: "bazinga_fail_12"}
	var message generics.Message = ""
	var positive bool
	ToggleHidePostTag(&disc, "uuid_1", &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, DocumentNotFoundError+"\n", message)
	disc.UUID = "testPost_forTags"
	message = ""
	positive = false
	ToggleHidePostTag(&disc, "trolololol", &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, CouldNotFindPostTag+"\n", message)
	message = ""
	positive = false
	ToggleHidePostTag(&disc, "uuid_1", &message, &positive)
	assert.Equal(t, true, positive)
	assert.Equal(t, SuccessfulPostTagHidden+"\n", message)
	message = ""
	positive = false
	ToggleHidePostTag(&disc, "uuid_1", &message, &positive)
	assert.Equal(t, true, positive)
	assert.Equal(t, SuccessfulPostTagUnhidden+"\n", message)
}

func testAddTag(t *testing.T) {
	tags := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A"}
	for _, tag := range tags {
		tag := tag
		go func() {
			var message generics.Message = ""
			var positive bool
			post := database.Document{UUID: "testPost_forTags"}

			AddPostTag(&post, &database.Posts{
				UUID:      "uuid_" + tag,
				Hidden:    false,
				Submitted: time.Now().UTC(),
				Info:      tag,
				Color:     "#00000" + tag,
			}, &message, &positive)
			assert.Equal(t, true, positive)
			assert.Equal(t, SuccessfulPostTag+"\n", message)
		}()
	}

	var message generics.Message = ""
	var positive bool
	disc := database.Document{UUID: "fail_request"}
	AddPostTag(&disc, &database.Posts{}, &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, DocumentNotFoundError+"\n", message)
	message = ""
	positive = false

	disc.UUID = "testDiscussion_forComments"
	AddPostTag(&disc, &database.Posts{}, &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, DocumentIsNotLegislativeText+"\n", message)

	time.Sleep(time.Millisecond * 100)
	documentLock.Lock()
	err := disc.GetByID("testPost_forTags")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(disc.Info.Post))
	err = disc.SaveChanges()
	assert.Nil(t, err)
	documentLock.Unlock()
}

func testHideComment(t *testing.T) {
	disc := database.Document{UUID: "bazinga_fail_12"}
	var message generics.Message = ""
	var positive bool
	ToggleHideDiscussionComment(&disc, "uuid_1", &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, DiscussioNotFoundError+"\n", message)
	disc.UUID = "testDiscussion_forComments"
	message = ""
	positive = false
	ToggleHideDiscussionComment(&disc, "trolololol", &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, CommentNotFoundError+"\n", message)
	message = ""
	positive = false
	ToggleHideDiscussionComment(&disc, "uuid_1", &message, &positive)
	assert.Equal(t, true, positive)
	assert.Equal(t, CommentSuccessfulHidden+"\n", message)
	message = ""
	positive = false
	ToggleHideDiscussionComment(&disc, "uuid_1", &message, &positive)
	assert.Equal(t, true, positive)
	assert.Equal(t, CommentSuccesfulUnhidden+"\n", message)
}

func testAddComment(t *testing.T) {
	comments := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, comment := range comments {
		comment := comment
		go func() {
			var message generics.Message = ""
			var positive bool
			disc := database.Document{UUID: "testDiscussion_forComments"}
			AddComment(&disc, &database.Discussions{
				UUID:        "uuid_" + comment,
				Hidden:      false,
				Written:     time.Now().UTC(),
				Author:      "test",
				Flair:       "",
				HTMLContent: comment,
			}, &message, &positive)
			assert.Equal(t, true, positive)
			assert.Equal(t, SuccessfulComment+"\n", message)
		}()
	}
	var message generics.Message = ""
	var positive bool
	disc := database.Document{UUID: "fail_request"}
	AddComment(&disc, &database.Discussions{}, &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, DiscussioNotFoundError+"\n", message)
	message = ""
	positive = false

	time.Sleep(time.Millisecond * 100)
	documentLock.Lock()
	err := disc.GetByID("testDiscussion_forComments")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(disc.Info.Discussion))
	disc.Info.Finishing = time.Now().UTC().Add(time.Hour * -24)
	err = disc.SaveChanges()
	assert.Nil(t, err)
	documentLock.Unlock()

	AddComment(&disc, &database.Discussions{}, &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, DiscussionAlreadyClosed+"\n", message)
	message = ""
	positive = false

	time.Sleep(1000)
	AddComment(&disc, &database.Discussions{}, &message, &positive)
	assert.Equal(t, false, positive)
	assert.Equal(t, CanNotCommentOnClosedDiscussion+"\n", message)
	message = ""
	positive = false
}

func testSetupDiscussionAndPost(t *testing.T) {
	doc := database.Document{
		UUID:         "testDiscussion_forComments",
		Written:      time.Time{},
		Organisation: "test",
		Type:         database.RunningDiscussion,
		Author:       "test",
		Flair:        "tset",
		Title:        "tset",
		Subtitle:     sql.NullString{},
		HTMLContent:  "test",
		Private:      false,
		Blocked:      false,
		Info: database.DocumentInfo{
			AnyPosterAllowed:          false,
			OrganisationPosterAllowed: false,
			Finishing:                 time.Now().UTC().Add(time.Hour * 24),
			Post:                      nil,
			Discussion:                []database.Discussions{},
			Votes:                     nil,
		},
		Viewer:  []database.Account{},
		Poster:  []database.Account{},
		Allowed: []database.Account{},
	}
	err := doc.CreateMe()
	assert.Nil(t, err)

	doc = database.Document{
		UUID:         "testPost_forTags",
		Written:      time.Time{},
		Organisation: "test",
		Type:         database.LegislativeText,
		Author:       "test",
		Flair:        "tset",
		Title:        "tset",
		Subtitle:     sql.NullString{},
		HTMLContent:  "test",
		Private:      false,
		Blocked:      false,
		Info: database.DocumentInfo{
			AnyPosterAllowed:          false,
			OrganisationPosterAllowed: false,
			Post:                      []database.Posts{},
			Discussion:                nil,
			Votes:                     nil,
		},
		Viewer:  []database.Account{},
		Poster:  []database.Account{},
		Allowed: []database.Account{},
	}
	err = doc.CreateMe()
	assert.Nil(t, err)
}
