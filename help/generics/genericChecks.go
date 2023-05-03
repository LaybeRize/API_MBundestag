package generics

import (
	"API_MBundestag/database"
	"fmt"
	"reflect"
)

func CheckTitelAndContentEmptyLayer[T, B any](message *T, content *B) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("Title").String() == "" ||
		ref.FieldByName("Content").String() == "" {
		//get the layer for the error and write back the message
		ref = reflect.ValueOf(message).Elem()
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(ContentAndTitelAreEmpty + "\n" + mesg)
		return true
	}
	return false
}

func CheckOrgOrTitle[T, B any](message *T, content *B) bool {
	ref := reflect.ValueOf(content).Elem()
	if ref.FieldByName("MainGroup").String() == "" ||
		ref.FieldByName("SubGroup").String() == "" ||
		ref.FieldByName("Name").String() == "" {
		ref = reflect.ValueOf(message).Elem()
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(NoMainGroupSubGroupOrNameProvided + "\n" + mesg)
		return true
	}
	return false
}

func CheckFieldNotEmpty[T any](v *T, fieldName string, errorMessage string) bool {
	ref := reflect.ValueOf(v).Elem()
	if ref.FieldByName(fieldName).String() == "" {
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(errorMessage + "\n" + mesg)
		return true
	}
	return false
}

func CheckOrgStatus[T any](message *T, statusString database.StatusString) bool {
	for _, r := range database.Stati {
		if r == string(statusString) {
			return false
		}
	}
	ref := reflect.ValueOf(message).Elem()
	mesg := ref.FieldByName("Message").String()
	ref.FieldByName("Message").SetString(StatusIsInvalid + "\n" + mesg)
	return true
}

func CheckLengthField[T any](v *T, length int, fieldName string, errorMessage string) bool {
	return checkLength(v, v, fieldName, length, errorMessage)
}

func CheckLengthFieldLayer[T, B any](message *T, content *B, length int, fieldName string, errorMessage string) bool {
	return checkLength(message, content, fieldName, length, errorMessage)
}

func CheckTitelAndContentEmpty[T any](v *T) bool {
	return CheckTitelAndContentEmptyLayer(v, v)
}

func CheckLengthContentLayer[T, B any](message *T, content *B, length int) bool {
	return checkLength(message, content, "Content", length, ContentTooLong)
}

func CheckLengthContent[T any](v *T, length int) bool {
	return checkLength(v, v, "Content", length, ContentTooLong)
}

func CheckLengthTitleLayer[T, B any](message *T, content *B, length int) bool {
	return checkLength(message, content, "Title", length, TitleTooLong)
}

func CheckLengthTitle[T any](v *T, length int) bool {
	return checkLength(v, v, "Title", length, TitleTooLong)
}

// CheckLengthSubtitleLayer checks if the Subtitle.String of the content exceeds the length of the parameter length
func CheckLengthSubtitleLayer[T, B any](message *T, content *B, length int) bool {
	ref := reflect.ValueOf(content).Elem()
	ref = ref.FieldByName("Subtitle")
	if len([]rune(ref.FieldByName("String").String())) > length {
		ref = reflect.ValueOf(message).Elem()
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(fmt.Sprintf(SubtitleTooLong, length) + "\n" + mesg)
		return true
	}
	return false
}

func CheckLengthSubtitle[T any](v *T, length int) bool {
	return checkLength(v, v, "Subtitle", length, SubtitleTooLong)
}

func checkLength[T, B any](message *T, content *B, fieldName string, length int, errorMsg string) bool {
	ref := reflect.ValueOf(content).Elem()
	if len([]rune(ref.FieldByName(fieldName).String())) > length {
		ref = reflect.ValueOf(message).Elem()
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(fmt.Sprintf(errorMsg, length) + "\n" + mesg)
		return true
	}
	return false
}

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
