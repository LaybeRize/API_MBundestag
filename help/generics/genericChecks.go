package generics

import (
	"API_MBundestag/database"
	"reflect"
)

type Message string

type MessageStruct struct {
	Message Message
	Positiv bool
}

func (m *Message) CheckTitleAndContentEmpty(content *any) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("Title").String() == "" ||
		ref.FieldByName("Content").String() == "" {
		// get the layer for the error and write back the message
		*m = ContentAndTitelAreEmpty + "\n" + *m
		return true
	}
	return false
}

func (m *Message) CheckOrgOrTitle(content any) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("MainGroup").String() == "" ||
		ref.FieldByName("SubGroup").String() == "" ||
		ref.FieldByName("Name").String() == "" {
		*m = NoMainGroupSubGroupOrNameProvided + "\n" + *m
		return true
	}
	return false
}

func (m *Message) CheckFieldNotEmpty(content any, fieldName string, errorMessage Message) bool {
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

func (m *Message) CheckLengthField(content any, fieldName string, length int, errorMessage Message) bool {
	ref := reflect.ValueOf(content).Elem()
	if len([]rune(ref.FieldByName(fieldName).String())) > length {
		*m = errorMessage + "\n" + *m
		return true
	}
	return false
}

func (m *Message) CheckLengthContent(content any, length int) bool {
	return m.CheckLengthField(content, "Content", length, ContentTooLong)
}

func (m *Message) CheckLengthTitle(content any, length int) bool {
	return m.CheckLengthField(content, "Title", length, TitleTooLong)
}

// CheckLengthSubtitle checks if the Subtitle of the content exceeds the length of the parameter length
func (m *Message) CheckLengthSubtitle(content any, length int) bool {
	return m.CheckLengthField(content, "Subtitle", length, SubtitleTooLong)
}

func (m *Message) CheckWriter(content any, writer *database.Account, acc *database.Account) bool {
	ref := reflect.ValueOf(content).Elem()
	err := writer.GetByDisplayNameWithParent(ref.FieldByName("SelectedAccount").String())
	if err != nil {
		*m = AccountDoesNotExists + "\n" + *m
		return true
	}

	if writer.Suspended ||
		(writer.Parent == nil && writer.ID != acc.ID) ||
		(writer.Parent != nil && writer.Parent.ID != acc.ID) {
		*m = AccountIsNotYours + "\n" + *m
		return true
	}
	return false
}

func (m *Message) CheckOrgExists(content any, org *database.Organisation) bool {
	ref := reflect.ValueOf(content).Elem()
	err := org.GetByName(ref.FieldByName("SelectedOrganisation").String())
	if err != nil {
		*m = OrganisationDoesNotExist + "\n" + *m
		return true
	}
	return false
}
