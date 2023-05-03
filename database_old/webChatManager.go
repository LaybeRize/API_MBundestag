package database

var ChatSchema = `
CREATE TABLE IF NOT EXISTS chats (
    uuid text UNIQUE NOT NULL,
    user1 TEXT NOT NULL,
    user2 TEXT NOT NULL,
    username1 TEXT NOT NULL,
    username2 TEXT NOT NULL,
    message1 TEXT NOT NULL,
    message2 TEXT NOT NULL
);
`

type Chat struct {
	UUID      string
	User1     string
	User2     string
	Username1 string
	Username2 string
	Message1  string
	Message2  string
}

func TestChatsDB() {
	TestDatabase("DROP TABLE IF EXISTS chats;", "")
	InitChatsDatabase()
}

func InitChatsDatabase() {
	DB.MustExec(ChatSchema)
}

func (chat *Chat) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO chats (uuid, user1, user2, username1, username2, message1, message2) VALUES (:uuid, :user1, :user2, :username1, :username2, :message1, :message2)", map[string]interface{}{
		"uuid":      chat.UUID,
		"user1":     chat.User1,
		"user2":     chat.User2,
		"username1": chat.Username1,
		"username2": chat.Username2,
		"message1":  chat.Message1,
		"message2":  chat.Message2,
	})
	return
}

func (chat *Chat) GetByID(uuid string, userUUID string) (err error) {
	err = DB.Get(chat, "SELECT * FROM chats WHERE uuid=$1 AND (user1 = $2 OR user2 = $2);", uuid, userUUID)
	return
}

func (chat *Chat) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE chats SET user1=:user1, user2=:user2, message1=:message1, message2=:message2 WHERE uuid=:uuid", map[string]interface{}{
		"uuid":     chat.UUID,
		"user1":    chat.User1,
		"user2":    chat.User2,
		"message1": chat.Message1,
		"message2": chat.Message2,
	})
	return
}
