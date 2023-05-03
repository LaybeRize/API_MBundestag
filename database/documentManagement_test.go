package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

var docAcc Account
var docAcc2 Account
var docAcc3 Account

func TestDocumentManagement(t *testing.T) {
	TestSetup()
	t.Run("testSetupForAccountsAndOrganisations", testSetupForAccountsAndOrganisations)
	t.Run("testCreateDocument", testCreateDocument)
	t.Run("testChangeAndGetDocument", testChangeAndGetDocument)
	t.Run("testDocumentCreate", testDocumentCreate)
	t.Run("testDocumentListUser", testDocumentListUser)
	t.Run("testDocumentListAdmin", testDocumentListAdmin)
}

func testDocumentListAdmin(t *testing.T) {
	list := DocumentList{}
	err, exists := list.GetAdminDocumentsAfter("test_doc_middle", 20, map[string]string{"blocked": "true"})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 0, len(list))
	err, exists = list.GetAdminDocumentsBefore("test_doc_middle", 20, map[string]string{"organisation": "test_docOrg_1", "blocked": "all"})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "test_doc_9", list[0].UUID)
	assert.Equal(t, "test_doc_8", list[1].UUID)
	assert.Equal(t, "test_doc_7", list[2].UUID)
	assert.Equal(t, "test_doc_1", list[3].UUID)
}

func testDocumentListUser(t *testing.T) {
	list := DocumentList{}
	err, exists := list.GetDocumentsAfter("test_doc_middle", 20, docAcc.ID, map[string]string{})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "test_doc_2", list[0].UUID)
	assert.Equal(t, "test_doc_4", list[1].UUID)
	assert.Equal(t, "test_doc_5", list[2].UUID)
	assert.Equal(t, "test_doc_6", list[3].UUID)
	err, exists = list.GetDocumentsBefore("test_doc_middle", 20, docAcc.ID, map[string]string{"organisation": "test_docOrg_1"})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "test_doc_8", list[0].UUID)
	assert.Equal(t, "test_doc_7", list[1].UUID)
	assert.Equal(t, "test_doc_1", list[2].UUID)
	err, exists = list.GetDocumentsBefore("test_doc_middle", 2, docAcc.ID, map[string]string{"organisation": "test_docOrg_1"})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test_doc_7", list[0].UUID)
	assert.Equal(t, "test_doc_1", list[1].UUID)

	err, exists = list.GetDocumentsAfter("", 20, docAcc3.ID, map[string]string{"organisation": "test_docOrg_1"})
	assert.Nil(t, err)
	assert.False(t, exists)
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "test_doc_8", list[0].UUID)
	assert.Equal(t, "test_doc_7", list[1].UUID)
	assert.Equal(t, "test_doc_4", list[2].UUID)
	assert.Equal(t, "test_doc_5", list[3].UUID)
	err, exists = list.GetDocumentsBefore("test_doc_middle", 20, docAcc3.ID, map[string]string{"organisation": "test_docOrg_1"})
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test_doc_8", list[0].UUID)
	assert.Equal(t, "test_doc_7", list[1].UUID)
}

func testDocumentCreate(t *testing.T) {
	doc := Document{
		UUID:         "test_doc_4",
		Written:      testGetTime(t, "1900-01-18"),
		Organisation: "test_docOrg_1",
		Type:         LegislativeText,
		Author:       "test",
		Flair:        "test",
		Title:        "test",
		HTMLContent:  "test",
	}
	err := doc.CreateMe()
	assert.Nil(t, err)
	doc.UUID, doc.Written = "test_doc_5", testGetTime(t, "1900-01-17")
	err = doc.CreateMe()
	assert.Nil(t, err)
	doc.UUID, doc.Written = "test_doc_6", testGetTime(t, "1900-01-16")
	doc.Private = true
	err = doc.CreateMe()
	assert.Nil(t, err)
	doc.Private = false
	doc.UUID, doc.Written = "test_doc_7", testGetTime(t, "1900-01-22")
	err = doc.CreateMe()
	assert.Nil(t, err)
	doc.UUID, doc.Written = "test_doc_8", testGetTime(t, "1900-01-23")
	doc.Type = FinishedDiscussion
	err = doc.CreateMe()
	assert.Nil(t, err)
	doc.UUID, doc.Written = "test_doc_9", testGetTime(t, "1900-01-24")
	doc.Blocked = true
	err = doc.CreateMe()
	assert.Nil(t, err)
}

func testChangeAndGetDocument(t *testing.T) {
	doc := Document{}
	err := doc.GetByIDOnlyWithAccount("test_doc_1", docAcc.ID)
	assert.Nil(t, err)
	err = doc.GetByIDOnlyWithAccount("test_doc_1", docAcc2.ID)
	assert.Nil(t, err)
	err = doc.GetByIDOnlyWithAccount("test_doc_1", docAcc3.ID)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	err = doc.GetByIDOnlyWithAccount("test_doc_2", docAcc.ID)
	assert.Nil(t, err)
	err = doc.GetByIDOnlyWithAccount("test_doc_2", docAcc2.ID)
	assert.Nil(t, err)
	err = doc.GetByIDOnlyWithAccount("test_doc_2", docAcc3.ID)
	assert.Nil(t, err)
}

func testCreateDocument(t *testing.T) {
	doc := Document{
		UUID:         "test_doc_1",
		Written:      testGetTime(t, "1900-01-21"),
		Organisation: "test_docOrg_1",
		Type:         LegislativeText,
		Author:       "test",
		Flair:        "test",
		Title:        "test",
		Subtitle:     sql.NullString{},
		HTMLContent:  "test",
		Private:      true,
		Blocked:      false,
		Info: DocumentInfo{
			AnyPosterAllowed: true,
			Post: []Posts{{
				UUID:      "test",
				Hidden:    true,
				Submitted: time.Time{},
				Info:      "asdasd",
				Color:     "#aaaaaa",
			}},
		},
		Viewer:  []Account{},
		Poster:  []Account{},
		Allowed: []Account{docAcc2},
	}
	err := doc.CreateMe()
	assert.Nil(t, err)
	second := Document{}
	err = second.GetByID("test_doc_1")
	assert.Nil(t, err)
	second.Written = second.Written.UTC()
	assert.Equal(t, doc, second)
	doc.Title = "test2"
	err = doc.SaveChanges()
	assert.Nil(t, err)
	assert.NotEqual(t, doc, second)
	err = second.GetByID("test_doc_1")
	assert.Nil(t, err)
	second.Written = second.Written.UTC()
	assert.Equal(t, doc, second)

	doc = Document{
		UUID:         "test_doc_2",
		Written:      testGetTime(t, "1900-01-19"),
		Organisation: "test_docOrg_2",
		Type:         LegislativeText,
		Author:       "test",
		Flair:        "test",
		Title:        "test",
		Subtitle:     sql.NullString{},
		HTMLContent:  "test",
		Private:      false,
		Blocked:      false,
		Info:         DocumentInfo{},
		Viewer:       []Account{},
		Poster:       []Account{},
		Allowed:      []Account{},
	}
	err = doc.CreateMe()
	assert.Nil(t, err)

	doc = Document{
		UUID:         "test_doc_middle",
		Written:      testGetTime(t, "1900-01-20"),
		Organisation: "test_docOrg_2",
		Type:         LegislativeText,
		Author:       "test",
		Flair:        "test",
		Title:        "test",
		HTMLContent:  "test",
	}
	err = doc.CreateMe()
	assert.Nil(t, err)
}

func testSetupForAccountsAndOrganisations(t *testing.T) {
	forTestCreateAccount(t, "doc_test_1", Account{Role: User})
	forTestCreateAccount(t, "doc_test_2", Account{Role: User})
	forTestCreateAccount(t, "doc_test_3", Account{Role: User})
	docAcc = Account{}
	docAcc2 = Account{}
	docAcc3 = Account{}
	err := docAcc.GetByUserName("doc_test_1")
	assert.Nil(t, err)
	err = docAcc2.GetByUserName("doc_test_2")
	assert.Nil(t, err)
	err = docAcc3.GetByUserName("doc_test_3")
	assert.Nil(t, err)
	org := Organisation{
		Name:      "test_docOrg_1",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{},
		Status:    Private,
		Members:   []Account{},
		Admins:    []Account{},
		Accounts:  []Account{docAcc},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = Organisation{
		Name:      "test_docOrg_2",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{},
		Status:    Private,
		Members:   []Account{},
		Admins:    []Account{},
		Accounts:  []Account{docAcc2},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
	org = Organisation{
		Name:      "test_docOrg_empty",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{},
		Status:    Private,
		Members:   []Account{},
		Admins:    []Account{},
		Accounts:  []Account{},
	}
	err = org.CreateMe()
	assert.Nil(t, err)
}
