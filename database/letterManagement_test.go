package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLetterManagment(t *testing.T) {
	TestSetup()
	t.Run("testCreateLetter", testCreateLetter)
	t.Run("testGetListOfLetters", testGetListOfLetters)
	t.Run("testGetModMailList", testGetModMailList)
	t.Run("testChanges", testChanges)
	t.Run("testCreateLists", testCreateLists)
}

var letterUserID = int64(0)

func testCreateLetter(t *testing.T) {
	acc := Account{}
	err := acc.GetByID(1)
	assert.Nil(t, err)
	letter := Letter{
		UUID:        "letter_test",
		Author:      acc.DisplayName,
		Flair:       acc.Flair,
		Title:       "letter_test",
		Content:     "letter_test",
		HTMLContent: "letter_test",
		ModMessage:  true,
		Info: LetterInfo{
			AllHaveToAgree:     false,
			NoSigning:          false,
			PeopleNotYetSigned: []string{},
			Signed:             []string{acc.DisplayName},
			Rejected:           []string{},
		},
		Viewer: []Account{acc},
	}
	err = letter.CreateMe()
	assert.Nil(t, err)
	forTestCreateAccount(t, "letter_user", Account{Role: User})
	err = acc.GetByUserName("letter_user")
	assert.Nil(t, err)
	letterUserID = acc.ID
	letter = Letter{
		UUID:        "letter_test2",
		Author:      acc.DisplayName,
		Flair:       acc.Flair,
		Title:       "letter_test2",
		Content:     "letter_test2",
		HTMLContent: "letter_test2",
		ModMessage:  true,
		Info: LetterInfo{
			AllHaveToAgree:     false,
			NoSigning:          false,
			PeopleNotYetSigned: []string{},
			Signed:             []string{acc.DisplayName},
			Rejected:           []string{},
		},
		Viewer: []Account{acc, acc},
	}
	err = letter.CreateMe()
	assert.Nil(t, err)
}

func testGetListOfLetters(t *testing.T) {
	list := LetterList{}
	getBasicLetterQuery("", 10, letterUserID).Find(&list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "letter_test2", list[0].UUID)
}

func testGetModMailList(t *testing.T) {
	list := LetterList{}
	getBasicModmailQuery("", 50).Find(&list)
	counter := 0
	for _, letter := range list {
		switch letter.UUID {
		case "letter_test2":
			fallthrough
		case "letter_test":
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}

func testChanges(t *testing.T) {
	letter := Letter{}
	err := letter.GetByID("letter_test2")
	assert.Nil(t, err)
	letter.Title = "test_change"
	letter.HTMLContent = "ajskdlajsldasd"
	err = letter.SaveChanges()
	assert.Nil(t, err)
	changed := Letter{}
	err = changed.GetByID("letter_test2")
	assert.Nil(t, err)
	assert.Equal(t, letter, changed)
	assert.Equal(t, 0, len(letter.Viewer))
	err = letter.GetByIDWithViewer("letter_test2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(letter.Viewer))
	assert.Equal(t, "letter_user", letter.Viewer[0].DisplayName)
}

func testCreateLists(t *testing.T) {
	forTestCreateAccount(t, "letter_user_special", Account{Role: User})
	acc := Account{}
	err := acc.GetByDisplayName("letter_user_special")
	assert.Nil(t, err)
	letterUserID = acc.ID
	temp := acc.ID
	testCreateSpecialLetter(t, "letter_old1", "1900-01-20", true)
	testCreateSpecialLetter(t, "letter_old2", "1900-01-19", true)
	testCreateSpecialLetter(t, "letter_old3", "1900-01-18", true)
	testCreateSpecialLetter(t, "letter_old4", "1900-01-17", true)
	testCreateSpecialLetter(t, "letter_old5", "1900-01-16", true)
	err = acc.GetByDisplayName("head_admin")
	assert.Nil(t, err)
	letterUserID = acc.ID
	testCreateSpecialLetter(t, "letter_old6", "1900-01-15", false)
	testCreateSpecialLetter(t, "letter_old7", "1900-01-14", false)
	list := LetterList{}
	var exists bool
	err, exists = list.GetLettersAfter("letter_old3", 4, letterUserID)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "letter_old6", list[0].UUID)
	assert.Equal(t, "letter_old7", list[1].UUID)
	letterUserID = temp
	err, exists = list.GetLettersAfter("letter_old3", 4, letterUserID)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "letter_old4", list[0].UUID)
	assert.Equal(t, "letter_old5", list[1].UUID)
	err, exists = list.GetLettersBefore("letter_old3", 4, letterUserID)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "letter_old1", list[0].UUID)
	assert.Equal(t, "letter_old2", list[1].UUID)
	err, exists = list.GetModMailsAfter("letter_old3", 4)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "letter_old4", list[0].UUID)
	assert.Equal(t, "letter_old5", list[1].UUID)
	err, exists = list.GetModMailsBefore("letter_old3", 2)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "letter_old1", list[0].UUID)
	assert.Equal(t, "letter_old2", list[1].UUID)
}

func testCreateSpecialLetter(t *testing.T, uuid string, timeStr string, modMessage bool) {
	acc := Account{}
	err := acc.GetByID(letterUserID)
	assert.Nil(t, err)
	l := &Letter{
		UUID:        uuid,
		Author:      acc.DisplayName,
		Flair:       acc.Flair,
		Written:     testGetTime(t, timeStr),
		Title:       uuid,
		Content:     uuid,
		HTMLContent: uuid,
		ModMessage:  modMessage,
		Info: LetterInfo{
			AllHaveToAgree:     false,
			NoSigning:          false,
			PeopleNotYetSigned: []string{},
			Signed:             []string{acc.DisplayName},
			Rejected:           []string{},
		},
		Viewer: []Account{acc},
	}
	err = db.Create(l).Error
	assert.Nil(t, err)
}

func testGetTime(t *testing.T, timeStr string) time.Time {
	ti, err := time.ParseInLocation("2006-01-02", timeStr, time.UTC)
	assert.Nil(t, err)
	return ti
}
