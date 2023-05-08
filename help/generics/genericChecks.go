package generics

import (
	"reflect"

	database "API_MBundestag/database"
)

type Message string

type MessageStruct struct {
	Message Message
	Positiv bool
}

var ContentOrTitleAreEmpty Message = "Inhalt oder Titel sind leer"

func (m *Message) CheckTitleAndContentEmpty(content *any) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("Title").String() == "" ||
		ref.FieldByName("Content").String() == "" {
		// get the layer for the error and write back the message
		*m = ContentOrTitleAreEmpty + "\n" + *m
		return true
	}
	return false
}

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

func (m *Message) CheckFieldNotEmpty(content *any, fieldName string, errorMessage Message) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName(fieldName).String() == "" {
		*m = errorMessage + "\n" + *m
		return true
	}
	return false
}

func (m *Message) CheckOrgStatus(statusString database.StatusString) bool {
	for _, r := range database.Stati {
		if r == string(statusString) {
			return false
		}
	}
	*m = StatusIsInvalid + "\n" + *m
	return true
}

func (m *Message) CheckLengthField(content *any, fieldName string, length int, errorMessage Message) bool {
	ref := reflect.ValueOf(content).Elem()
	if len([]rune(ref.FieldByName(fieldName).String())) > length {
		*m = errorMessage + "\n" + *m
		return true
	}
	return false
}

func (m *Message) CheckLengthContent(content *any, length int) bool {
	return m.CheckLengthField(content, "Content", length, ContentTooLong)
}

func (m *Message) CheckLengthTitle(content *any, length int) bool {
	return m.CheckLengthField(content, "Title", length, TitleTooLong)
}

// CheckLengthSubtitleLayer checks if the Subtitle.String of the content exceeds the length of the parameter length
func (m *Message) CheckLengthSubtitle(content *any, length int) bool {
	return m.CheckLengthField(content, "Subtitle", length, SubtitleTooLong)
}

/*

func CheckWriter[T any](v *T, writer *database.Account, acc *database.Account) bool {
	ref := reflect.ValueOf(v).Elem()
	mesg := ref.FieldByName("Message").String()
	err := writer.GetByDisplayName(ref.FieldByName("SelectedAccount").String())
	if err != nil {
		ref.FieldByName("Message").SetString(AccountDoesNotExists + "\n" + mesg)
		return true
	}
	if (writer.Linked.Int64 != acc.ID || writer.Suspended) && !(writer.DisplayName == acc.DisplayName) {
		ref.FieldByName("Message").SetString(AccountIsNotYours + "\n" + mesg)
		return true
	}
	return false
}

func CheckOrgExists[T any](v *T, org *database.Organisation) bool {
	ref := reflect.ValueOf(v).Elem()
	err := org.GetByName(ref.FieldByName("SelectedOrganisation").String())
	if err != nil {
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(OrganisationDoesNotExist + "\n" + mesg)
		return true
	}
	return false
}


*/
