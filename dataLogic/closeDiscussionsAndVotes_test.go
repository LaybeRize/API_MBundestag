package dataLogic

import (
	"API_MBundestag/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloseDiscussionOrVote(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupDiscussionsAndVotes", testSetupDiscussionsAndVotes)
	t.Run("testCloseDiscussion", testCloseDiscussion)
	t.Run("testCloseVote", testCloseVote)
}

func testCloseVote(t *testing.T) {
	err := CloseDiscussionOrVote("test_close_vote")
	assert.Nil(t, err)
	doc := database.Document{}
	err = doc.GetByID("test_close_vote")
	assert.Nil(t, err)
	assert.Equal(t, doc.Type, database.FinishedVote)
	vote := database.Votes{}
	err = vote.GetByID("test_close_single_vote")
	assert.Nil(t, err)
	assert.Equal(t, true, vote.Finished)
	err = vote.GetByID("test_close_multiple_votes")
	assert.Nil(t, err)
	assert.Equal(t, true, vote.Finished)
}

func testCloseDiscussion(t *testing.T) {
	err := CloseDiscussionOrVote("test_close_discussion")
	assert.Nil(t, err)
	doc := database.Document{}
	err = doc.GetByID("test_close_discussion")
	assert.Nil(t, err)
	assert.Equal(t, doc.Type, database.FinishedDiscussion)
}

func testSetupDiscussionsAndVotes(t *testing.T) {
	doc := database.Document{
		UUID:         "test_close_discussion",
		Organisation: "empty",
		Type:         database.RunningDiscussion,
		Author:       "test",
		Flair:        "test",
		Title:        "test",
		HTMLContent:  "test",
		Viewer:       []database.Account{},
		Poster:       []database.Account{},
		Allowed:      []database.Account{},
	}
	err := doc.CreateMe()
	assert.Nil(t, err)

	doc = database.Document{
		UUID:         "test_close_vote",
		Organisation: "empty",
		Type:         database.RunningVote,
		Author:       "XXXX",
		Flair:        "XXXX",
		Title:        "XXXX",
		HTMLContent:  "XXXX",
		Viewer:       []database.Account{},
		Poster:       []database.Account{},
		Allowed:      []database.Account{},
		Info:         database.DocumentInfo{Votes: []string{"test_close_single_vote", "test_close_multiple_votes"}},
	}
	err = doc.CreateMe()
	assert.Nil(t, err)
	vote := database.Votes{
		UUID:                   "test_close_single_vote",
		Parent:                 "test_close_vote",
		Question:               "Test Question",
		ShowNumbersWhileVoting: false,
		ShowNamesWhileVoting:   false,
		ShowNamesAfterVoting:   false,
		Finished:               false,
		Info:                   database.VoteInfo{},
	}
	err = vote.CreateMe()
	assert.Nil(t, err)
	vote = database.Votes{
		UUID:                   "test_close_multiple_votes",
		Parent:                 "test_close_vote",
		Question:               "Test Question",
		ShowNumbersWhileVoting: false,
		ShowNamesWhileVoting:   false,
		ShowNamesAfterVoting:   false,
		Finished:               false,
		Info:                   database.VoteInfo{},
	}
	err = vote.CreateMe()
	assert.Nil(t, err)
}
