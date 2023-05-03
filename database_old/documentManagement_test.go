package database

import (
	"API_MBundestag/help"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var doc Document
var docList DocumentList

func TestDocuments(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestDocumentsDB()
	TestOrganisationsDB()

	t.Run("testCreateLetter", testCreateDocument)
	t.Run("testReadLetter", testReadDocument)
	t.Run("testChangeLetter", testChangeDocument)
	t.Run("createQueryiableDocuments", createQueryableDocuments)
	t.Run("testAfterDocument", testAfterDocument)
	t.Run("testBeforeDocument", testBeforeDocument)
}

func testBeforeDocument(t *testing.T) {
	res := DocumentList{}
	err, exists := res.GetGetDocumentsBefore(docList[7].UUID, 20, map[string]string{})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 7, len(res))
	assert.Equal(t, docList[0].UUID, res[0].UUID)
	assert.Equal(t, docList[2].UUID, res[2].UUID)
	assert.Equal(t, docList[6].UUID, res[6].UUID)
	err, exists = res.GetGetDocumentsBefore(docList[7].UUID, 2, map[string]string{})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, docList[5].UUID, res[0].UUID)
	assert.Equal(t, docList[6].UUID, res[1].UUID)
	err, exists = res.GetGetDocumentsBefore(docList[7].UUID, 20, map[string]string{"title": "3", "author": "test12"})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(res))
}

func testAfterDocument(t *testing.T) {
	res := DocumentList{}
	err, exists := res.GetDocumentsAfter(docList[0].UUID, 20, map[string]string{})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 7, len(res))
	assert.Equal(t, docList[1].UUID, res[0].UUID)
	assert.Equal(t, docList[3].UUID, res[2].UUID)
	assert.Equal(t, docList[7].UUID, res[6].UUID)
	err, exists = res.GetDocumentsAfter(docList[0].UUID, 2, map[string]string{})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, docList[1].UUID, res[0].UUID)
	assert.Equal(t, docList[2].UUID, res[1].UUID)
	err, exists = res.GetDocumentsAfter(docList[0].UUID, 20, map[string]string{"organisation": "test"})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, docList[4].UUID, res[0].UUID)
	assert.Equal(t, docList[7].UUID, res[3].UUID)
	err, exists = res.GetDocumentsAfter(docList[0].UUID, 20, map[string]string{"title": "e"})
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 7, len(res))
	assert.Equal(t, docList[1].UUID, res[0].UUID)
	assert.Equal(t, docList[3].UUID, res[2].UUID)
	assert.Equal(t, docList[7].UUID, res[6].UUID)
	err, exists = res.GetDocumentsAfter(docList[0].UUID, 20, map[string]string{"title": "e"}, LegislativeText)
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, docList[1].UUID, res[0].UUID)
	assert.Equal(t, docList[5].UUID, res[2].UUID)
	err, exists = res.GetDocumentsAfter(docList[0].UUID, 20, map[string]string{"title": "e"}, LegislativeText, FinishedVote)
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, docList[1].UUID, res[0].UUID)
	assert.Equal(t, docList[5].UUID, res[2].UUID)
	assert.Equal(t, docList[6].UUID, res[3].UUID)
}

func createQueryableDocuments(t *testing.T) {
	var err error
	docList = DocumentList{doc}
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Type = FinishedVote
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Type = LegislativeText
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Type = UnfinishedVote
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Organisation = "test2"
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Type = LegislativeText
	doc.Author = "test2"
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Type = LegislativeText
	doc.Title = "test3"
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)
	time.Sleep(100)
	doc.UUID = uuid.New().String()
	doc.Type = Discussion
	doc.Title = "end"
	err = doc.CreateMe()
	assert.Nil(t, err)
	docList = append(DocumentList{doc}, docList...)

	assert.Equal(t, 8, len(docList))
}

func testChangeDocument(t *testing.T) {
	doc.Info.Viewer = []string{"sdasd", "vcxvyxdf", "dsfad"}
	doc.Private = false
	err := doc.SaveChanges()
	assert.Nil(t, err)
	res := Document{}
	err = res.GetByID(doc.UUID)
	assert.Nil(t, err)
	res.Written = time.Time{}
	assert.Equal(t, doc, res)
}

func testReadDocument(t *testing.T) {
	res := Document{}
	doc.Written = time.Time{}
	err := res.GetByID(doc.UUID)
	assert.Nil(t, err)
	res.Written = time.Time{}
	assert.Equal(t, doc, res)
}

func testCreateDocument(t *testing.T) {
	doc = Document{
		UUID:         uuid.New().String(),
		Organisation: "test",
		Type:         Discussion,
		Author:       "test",
		Flair:        "test",
		Title:        "test",
		Subtitle:     sql.NullString{},
		HTMLContent:  "test",
		Private:      true,
		Info: DocumentInfo{
			Viewer: []string{"test", "test"},
			Poster: []string{"test"},
		},
	}
	err := doc.CreateMe()
	assert.Nil(t, err)
}

//TODO add tests for the variouse cases
