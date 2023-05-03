package help

import "reflect"

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

var NoMainGroupSubGroupOrNameProvided Message = "Es wurde keine Hauptkategorie, Unterkategorie oder ein Name angegeben"

func (m *Message) CheckOrgOrTitle(content *any) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("MainGroup").String() == "" ||
		ref.FieldByName("SubGroup").String() == "" ||
		ref.FieldByName("Name").String() == "" {
		*m = NoMainGroupSubGroupOrNameProvided + "\n" + *m
		return true
	}
	return false
}
