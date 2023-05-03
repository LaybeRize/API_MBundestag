package database

import (
	"API_MBundestag/help"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var letter Letter

func TestLetters(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestLettersDB()

	t.Run("testCreateLetter", testCreateLetter)
	t.Run("testReadLetter", testReadLetter)
	t.Run("testChangeLetter", testChangeLetter)
	t.Run("createNewLetters", createNewLetters)
	t.Run("testReadLetterList", testReadLetterList)
}

func testReadLetterList(t *testing.T) {
	list := LetterList{}
	err, ex := list.GetPublicationAfter("", 2, "2asd", false)
	assert.False(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test2", list[0].Title)
	assert.Equal(t, "werfy", list[1].Title)
	err, ex = list.GetPublicationBefore(list[1].UUID, 2, "", true)
	assert.True(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test3", list[0].Title)
	assert.Equal(t, "test2", list[1].Title)
	err, ex = list.GetPublicationBefore("", 50, "2asd", false)
	assert.False(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))
	err, ex = list.GetPublicationBefore("test2", 50, "", true)
	assert.True(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "test3", list[0].Title)
}

func createNewLetters(t *testing.T) {
	letter.UUID = "test2"
	letter.Title = "test2"
	err := letter.CreateMe()
	assert.Nil(t, err)

	letter.Info.PeopleInvitedToSign = []string{}
	letter.UUID = "test3"
	letter.Title = "test3"
	err = letter.CreateMe()
	assert.Nil(t, err)
}

func testChangeLetter(t *testing.T) {
	letter.Info.PeopleInvitedToSign = append(letter.Info.PeopleInvitedToSign, "asdv2")
	letter.Info.AllHaveToAgree = true
	err := letter.SaveChanges()
	assert.Nil(t, err)
	res := Letter{}
	err = res.GetByID("test")
	assert.Nil(t, err)
	res.Written = time.Time{}
	assert.Equal(t, letter, res)
}

func testReadLetter(t *testing.T) {
	res := Letter{}
	err := res.GetByID("test")
	assert.Nil(t, err)
	res.Written = time.Time{}
	assert.Equal(t, letter, res)
}

func testCreateLetter(t *testing.T) {
	letter = Letter{
		UUID:        "test",
		Author:      "vcbx",
		Written:     time.Time{},
		Title:       "werfy",
		HTMLContent: "asdfer",
		Info: LetterInfo{
			ModMessage:          true,
			AllHaveToAgree:      false,
			PeopleInvitedToSign: []string{"2asd", "sdaf"},
			Signed:              []string{"aysdas"},
			Rejected:            []string{"xycx", "asdydsc", "as"},
		},
	}
	err := letter.CreateMe()
	assert.Nil(t, err)
}
