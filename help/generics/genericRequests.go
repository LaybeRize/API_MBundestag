package generics

import (
	"reflect"
)

type Message string

type MessageStruct struct {
	Message Message
	Positiv bool
}

var ContentOrTitelAreEmpty Message = "Inhalt oder Titel sind leer"

func (m *Message) CheckTitelAndContentEmptyLayer(content *any) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("Title").String() == "" ||
		ref.FieldByName("Content").String() == "" {
		//get the layer for the error and write back the message
		*m = ContentOrTitelAreEmpty + "\n" + *m
		return true
	}
	return false
}
