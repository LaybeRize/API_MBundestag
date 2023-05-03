package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

func (msg *Messages) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &msg)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &msg)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (msg *Messages) Value() driver.Value {
	l, _ := json.Marshal(&msg)
	return l
}

type (
	Messages struct {
		User1 []string
		User2 []string
	}
	Chat struct {
		UUID      string `gorm:"primaryKey"`
		User1     string
		User2     string
		Username1 string
		Username2 string
		Messages  Messages `gorm:"type:jsonb"`
	}
)

func (chat *Chat) CreateMe() (err error) {
	err = db.Create(chat).Error
	return
}

func (chat *Chat) GetByID(uuid string) (err error) {
	*chat = Chat{}
	err = db.First(chat, "uuid = ?", uuid).Error
	return
}

func (chat *Chat) SaveChanges() (err error) {
	err = db.Save(chat).Error
	return
}
