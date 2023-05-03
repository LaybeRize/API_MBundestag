package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChatManagement(t *testing.T) {
	TestSetup()
	t.Run("testCreateChat", testCreateChat)
	t.Run("testChangeChat", testChangeChat)
}

func testChangeChat(t *testing.T) {
	chat := Chat{}
	err := chat.GetByID("chat_test1")
	assert.Nil(t, err)
	chat.User2 = "fasdbassd"
	chat.Messages.User1 = []string{"sdafasd", "bsdsad"}
	err = chat.SaveChanges()
	assert.Nil(t, err)
	second := Chat{}
	err = second.GetByID("chat_test1")
	assert.Nil(t, err)
	assert.Equal(t, chat, second)
}

func testCreateChat(t *testing.T) {
	chat := Chat{
		UUID:      "chat_test1",
		User1:     "asda",
		User2:     "basd",
		Username1: "qwe",
		Username2: "terfsd",
		Messages: Messages{
			User1: []string{},
			User2: []string{"test"},
		},
	}
	err := chat.CreateMe()
	assert.Nil(t, err)
	second := Chat{}
	err = second.GetByID("chat_test1")
	assert.Nil(t, err)
	assert.Equal(t, chat, second)
}
