package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoteManagement(t *testing.T) {
	TestSetup()
	t.Run("testCreateVote", testCreateVote)
	t.Run("testChangeVote", testChangeVote)
}

func testChangeVote(t *testing.T) {
	vote := Votes{}
	err := vote.GetByID("vote_test1")
	assert.Nil(t, err)
	vote.Question = "lol"
	vote.Finished = true
	vote.Info.MaxPosition = 20
	vote.Info.Options = []string{"asdnalsd", "a", "vasdasd"}
	err = vote.SaveChanges()
	assert.Nil(t, err)
	second := Votes{}
	err = second.GetByID("vote_test1")
	assert.Nil(t, err)
	assert.Equal(t, vote, second)
}

func testCreateVote(t *testing.T) {
	vote := Votes{
		UUID:                   "vote_test1",
		Parent:                 "test2",
		Question:               "bazinga",
		ShowNumbersWhileVoting: false,
		ShowNamesWhileVoting:   true,
		ShowNamesAfterVoting:   false,
		Finished:               false,
		Info: VoteInfo{
			Results: map[string]Results{"test": {
				Votee:       "bazinga",
				InvalidVote: false,
				Votes:       map[string]int{"a": 1},
			}},
			Summary:     Summary{},
			VoteMethod:  SingleVote,
			MaxPosition: 12,
			Options:     []string{"test", "a", "lol"},
		},
	}
	err := vote.CreateMe()
	assert.Nil(t, err)
	second := Votes{}
	err = second.GetByID("vote_test1")
	assert.Nil(t, err)
	assert.Equal(t, vote, second)
}
