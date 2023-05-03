package dataLogic

import (
	"API_MBundestag/database"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedTitleHierarchy = TitleMainGroup{}
var expectedTitleHierarchy2 = TitleMainGroup{}

func TestTitleHierarchy(t *testing.T) {
	database.TestSetup()

	t.Run("testSingleTitleHierarchy", testSingleTitleHierarchy)
	t.Run("testSingleTitlesInDifferentMainGroups", testSingleTitlesInDifferentMainGroups)
	t.Run("testMultipleTitlesInDifferentMainGroups", testMultipleTitlesInDifferentMainGroups)
}

func testMultipleTitlesInDifferentMainGroups(t *testing.T) {
	title := database.Title{
		Name:      "test_titleHierarchy3",
		MainGroup: "test_titleHierarchy2",
		SubGroup:  "test_titleHierarchy3",
		Flair: sql.NullString{
			String: "",
			Valid:  false,
		},
		Holder: []database.Account{},
	}
	err := title.CreateMe()
	assert.Nil(t, err)
	err = RefreshTitleHierarchy()
	assert.Nil(t, err)

	expectedTitleHierarchy2.Groups = append(expectedTitleHierarchy2.Groups, TitleSubGroup{
		Name:   "test_titleHierarchy3",
		Titles: []database.Title{title},
	})
	var ref TitleMainGroup
	var ref2 TitleMainGroup
	counter := 0
	for i, titleGroup := range GetTitleHierarchy() {
		if titleGroup.Name == "test_titleHierarchy" {
			ref = titleGroup
			counter = i
		}
		if titleGroup.Name == "test_titleHierarchy2" {
			ref2 = titleGroup
			if counter+1 != i {
				assert.Fail(t, "wrong order")
			}
		}
	}
	assert.Equal(t, expectedTitleHierarchy, ref)
	assert.Equal(t, expectedTitleHierarchy2, ref2)
}

func testSingleTitlesInDifferentMainGroups(t *testing.T) {
	title := database.Title{
		Name:      "test_titleHierarchy2",
		MainGroup: "test_titleHierarchy2",
		SubGroup:  "test_titleHierarchy2",
		Flair: sql.NullString{
			String: "",
			Valid:  false,
		},
		Holder: []database.Account{},
	}
	err := title.CreateMe()
	assert.Nil(t, err)
	err = RefreshTitleHierarchy()
	assert.Nil(t, err)

	expectedTitleHierarchy2 = TitleMainGroup{
		Name: "test_titleHierarchy2",
		Groups: []TitleSubGroup{{
			Name:   "test_titleHierarchy2",
			Titles: []database.Title{title},
		}},
	}
	var ref TitleMainGroup
	var ref2 TitleMainGroup
	counter := 0
	for i, titleGroup := range GetTitleHierarchy() {
		if titleGroup.Name == "test_titleHierarchy" {
			ref = titleGroup
			counter = i
		}
		if titleGroup.Name == "test_titleHierarchy2" {
			ref2 = titleGroup
			if counter+1 != i {
				assert.Fail(t, "wrong order")
			}
		}
	}
	assert.Equal(t, expectedTitleHierarchy, ref)
	assert.Equal(t, expectedTitleHierarchy2, ref2)
}

func testSingleTitleHierarchy(t *testing.T) {
	title := database.Title{
		Name:      "test_titleHierarchy",
		MainGroup: "test_titleHierarchy",
		SubGroup:  "test_titleHierarchy",
		Flair: sql.NullString{
			String: "",
			Valid:  false,
		},
		Holder: []database.Account{},
	}
	err := title.CreateMe()
	assert.Nil(t, err)
	err = RefreshTitleHierarchy()
	assert.Nil(t, err)

	expectedTitleHierarchy = TitleMainGroup{
		Name: "test_titleHierarchy",
		Groups: []TitleSubGroup{{
			Name:   "test_titleHierarchy",
			Titles: []database.Title{title},
		}},
	}
	var ref TitleMainGroup
	for _, titleGroup := range GetTitleHierarchy() {
		if titleGroup.Name == "test_titleHierarchy" {
			ref = titleGroup
		}
	}
	assert.Equal(t, expectedTitleHierarchy, ref)
}
