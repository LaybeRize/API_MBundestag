package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestVoteLogic(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupDocumentsAndVotes", testSetupDocumentsAndVotes)
	t.Run("testFailVoteFinish", testFailVoteFinish)
	t.Run("testFailDocument", testFailDocument)
	t.Run("testSucessVote", testSucessVote)
	t.Run("testFailAlreadyVoted", testFailAlreadyVoted)
}

func testFailAlreadyVoted(t *testing.T) {
	time.Sleep(time.Millisecond * 100)
	vote := database.Votes{UUID: "testVoteLogicSingle"}
	var msg help.Message
	var positive bool
	err := AddResultForUser(&vote, map[string]int{}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, YouAlreadyVoted+"\n", msg)
	assert.False(t, positive)

	msg = ""
	vote.UUID = "testVoteLogicMultiple"
	err = AddResultForUser(&vote, map[string]int{}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, YouAlreadyVoted+"\n", msg)
	assert.False(t, positive)

	msg = ""
	vote.UUID = "testVoteLogicThree"
	err = AddResultForUser(&vote, map[string]int{}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, YouAlreadyVoted+"\n", msg)
	assert.False(t, positive)

	msg = ""
	vote.UUID = "testVoteLogicRanked"
	err = AddResultForUser(&vote, map[string]int{}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, YouAlreadyVoted+"\n", msg)
	assert.False(t, positive)
}

func testSucessVote(t *testing.T) {
	vote := database.Votes{UUID: "testVoteLogicSingle"}
	var msg help.Message
	var positive bool
	err := AddResultForUser(&vote, map[string]int{
		"test 1": 0,
		"test 2": 1,
		"test 3": 0,
	}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, SuccessfulVote+"\n", msg)
	assert.True(t, positive)

	time.Sleep(time.Millisecond * 30)
	err = vote.GetByID(vote.UUID)
	assert.Nil(t, err)
	assert.Equal(t, map[string]database.Results{
		"test": {
			Votee:       "test",
			InvalidVote: false,
			Votes: map[string]int{
				"test 1": 0,
				"test 2": 1,
				"test 3": 0,
			},
		},
	}, vote.Info.Results)
	assert.Equal(t, map[string]int{
		"test 1": 0,
		"test 2": 1,
		"test 3": 0,
	}, vote.Info.Summary.Sums)
	assert.Equal(t, map[string]string{
		"test": "test 2",
	}, vote.Info.Summary.Person)

	positive = false
	msg = ""
	vote.UUID = "testVoteLogicMultiple"
	err = AddResultForUser(&vote, map[string]int{
		"test 1": 0,
		"test 2": 1,
		"test 3": 1,
	}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, SuccessfulVote+"\n", msg)
	assert.True(t, positive)

	time.Sleep(time.Millisecond * 30)
	err = vote.GetByID(vote.UUID)
	assert.Nil(t, err)
	assert.Equal(t, map[string]database.Results{
		"test": {
			Votee:       "test",
			InvalidVote: false,
			Votes: map[string]int{
				"test 1": 0,
				"test 2": 1,
				"test 3": 1,
			},
		},
	}, vote.Info.Results)
	assert.Equal(t, map[string]int{
		"test 1": 0,
		"test 2": 1,
		"test 3": 1,
	}, vote.Info.Summary.Sums)
	assert.Equal(t, map[string]string{
		"test": "test 2, test 3",
	}, vote.Info.Summary.Person)

	positive = false
	msg = ""
	vote.UUID = "testVoteLogicThree"
	err = AddResultForUser(&vote, map[string]int{
		"test 1": 0,
		"test 2": -1,
		"test 3": 1,
	}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, SuccessfulVote+"\n", msg)
	assert.True(t, positive)

	time.Sleep(time.Millisecond * 30)
	err = vote.GetByID(vote.UUID)
	assert.Nil(t, err)
	assert.Equal(t, map[string]database.Results{
		"test": {
			Votee:       "test",
			InvalidVote: false,
			Votes: map[string]int{
				"test 1": 0,
				"test 2": -1,
				"test 3": 1,
			},
		},
	}, vote.Info.Results)
	assert.Equal(t, map[string]int{
		"test 1": 0,
		"test 2": -1,
		"test 3": 1,
	}, vote.Info.Summary.Sums)
	assert.Equal(t, map[string]string{
		"test": ForVote + "test 3\n" + AgainstVote + "test 2",
	}, vote.Info.Summary.Person)

	positive = false
	msg = ""
	vote.UUID = "testVoteLogicRanked"
	err = AddResultForUser(&vote, map[string]int{
		"test 1": 0,
		"test 2": 2,
		"test 3": 1,
	}, false, "test", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, SuccessfulVote+"\n", msg)
	assert.True(t, positive)

	time.Sleep(time.Millisecond * 30)
	err = vote.GetByID(vote.UUID)
	assert.Nil(t, err)
	assert.Equal(t, map[string]database.Results{
		"test": {
			Votee:       "test",
			InvalidVote: false,
			Votes: map[string]int{
				"test 1": 0,
				"test 2": 2,
				"test 3": 1,
			},
		},
	}, vote.Info.Results)
	assert.Equal(t, map[string]map[string]int{
		"test": {
			"test 1": 0,
			"test 2": 2,
			"test 3": 1,
		},
	}, vote.Info.Summary.RankedMap)
}

func testFailDocument(t *testing.T) {
	vote := database.Votes{UUID: "testVoteLogicFail"}
	var msg help.Message
	var positive bool
	err := AddResultForUser(&vote, map[string]int{}, false, "", &msg, &positive)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, ErrorWhileSavingVote+"\n", msg)
	assert.False(t, positive)
	msg = ""

	vote.UUID = "testVoteLogicFail1"
	err = AddResultForUser(&vote, map[string]int{}, false, "", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, VoteAlreadyFinished+"\n", msg)
	assert.False(t, positive)
	msg = ""

	time.Sleep(time.Millisecond * 100)
	doc := database.Document{}
	err = doc.GetByID("testVoteLogicFail")
	assert.Equal(t, doc.Type, database.FinishedVote)
	err = vote.GetByID("testVoteLogicFail1")
	assert.True(t, vote.Finished)

	vote.UUID = "testVoteLogicFail2"
	err = AddResultForUser(&vote, map[string]int{}, false, "", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, ErrorNotAVote+"\n", msg)
	assert.False(t, positive)
	msg = ""
}

func testFailVoteFinish(t *testing.T) {
	vote := database.Votes{UUID: "alsdopawbejqnkie"}
	var msg help.Message
	var positive bool
	err := AddResultForUser(&vote, map[string]int{}, false, "", &msg, &positive)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, VoteAlreadyFinished+"\n", msg)
	assert.False(t, positive)
	msg = ""

	vote.UUID = "testVoteLogic"
	err = AddResultForUser(&vote, map[string]int{}, true, "", &msg, &positive)
	assert.Nil(t, err)
	assert.Equal(t, VoteAlreadyFinished+"\n", msg)
	assert.False(t, positive)
}

func testSetupDocumentsAndVotes(t *testing.T) {
	vote := database.Votes{
		UUID:     "testVoteLogic",
		Question: "TestQuestion",
		Finished: true,
	}
	err := vote.CreateMe()
	assert.Nil(t, err)
	vote = database.Votes{
		UUID:     "testVoteLogicFail",
		Question: "TestQuestion",
	}
	err = vote.CreateMe()
	assert.Nil(t, err)
	vote = database.Votes{
		UUID:     "testVoteLogicFail1",
		Parent:   "testVoteLogicFail",
		Question: "TestQuestion",
	}
	err = vote.CreateMe()
	assert.Nil(t, err)
	vote = database.Votes{
		UUID:     "testVoteLogicFail2",
		Parent:   "testVoteLogicFail1",
		Question: "TestQuestion",
	}
	err = vote.CreateMe()
	assert.Nil(t, err)
	doc := database.Document{
		UUID: "testVoteLogicFail",
		Type: database.RunningVote,
		Info: database.DocumentInfo{
			Finishing: time.Now().UTC().Add(time.Hour * -24),
			Votes:     []string{"testVoteLogicFail1"},
		},
	}
	err = doc.CreateMe()
	assert.Nil(t, err)
	doc = database.Document{
		UUID: "testVoteLogicFail1",
		Type: database.LegislativeText,
	}
	err = doc.CreateMe()
	assert.Nil(t, err)
	vote = database.Votes{
		UUID:     "testVoteLogicSingle",
		Parent:   "testVoteLogicSuccess",
		Question: "Test Question",
		Info: database.VoteInfo{
			Results: map[string]database.Results{},
			Summary: database.Summary{
				Sums:         map[string]int{},
				RankedMap:    map[string]map[string]int{},
				Person:       map[string]string{},
				InvalidVotes: []string{},
				CSV:          "",
			},
			VoteMethod:  database.SingleVote,
			MaxPosition: 10,
			Options:     []string{"test 1", "test 2", "test 3"},
		},
	}
	err = vote.CreateMe()
	assert.Nil(t, err)
	vote.UUID, vote.Info.VoteMethod = "testVoteLogicMultiple", database.MultipleVotes
	err = vote.CreateMe()
	assert.Nil(t, err)
	vote.UUID, vote.Info.VoteMethod = "testVoteLogicThree", database.ThreeCategoryVoting
	err = vote.CreateMe()
	assert.Nil(t, err)
	vote.UUID, vote.Info.VoteMethod = "testVoteLogicRanked", database.VoteRanking
	err = vote.CreateMe()
	assert.Nil(t, err)

	doc = database.Document{
		UUID: "testVoteLogicSuccess",
		Type: database.RunningVote,
		Info: database.DocumentInfo{
			Finishing: time.Now().UTC().Add(time.Hour * 24),
			Votes:     []string{"testVoteLogicSingle", "testVoteLogicMultiple", "testVoteLogicThree", "testVoteLogicRanked"},
		},
	}
	err = doc.CreateMe()
	assert.Nil(t, err)
}
