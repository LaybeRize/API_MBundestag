package help

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	gen "API_MBundestag/help/generics"
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

func FillAllNotSuspendedNames[T any](v *T) {
	names, err := dataLogic.GetAllAccountNamesNotSuspended()
	slice := reflect.ValueOf(names)
	updateField(v, "Names", slice, err, gen.NamesQueryError)
}

func FillUserAndDisplayNames[T any](v *T) {
	names := database.NameList{}
	err := names.GetAllUserAndDisplayName()
	slice := reflect.ValueOf(names)
	updateField(v, "Names", slice, err, gen.NamesQueryError)
}

func FillOrganisationNames[T any](v *T) {
	orgNames, err := dataLogic.GetAllOrganisationNames()
	slice := reflect.ValueOf(orgNames)
	updateField(v, "OrgNames", slice, err, gen.OrgNamesQueryError)
}

func FillOrganisationGroups[T any](v *T) {
	main, sub, err := dataLogic.GetNamesForSubAndMainGroups()
	sliceMain := reflect.ValueOf(main)
	sliceSub := reflect.ValueOf(sub)
	updateField(v, "ExistingMainGroup", sliceMain, nil, "")
	updateField(v, "ExistingSubGroup", sliceSub, err, gen.GroupQueryError)
}

func FillOwnAccounts[T any](v *T, acc *database.Account) {
	ownAccounts := database.AccountList{}
	err := ownAccounts.GetAllPressAccountsFromAccountPlusSelf(acc)
	slice := reflect.ValueOf(ownAccounts)
	updateField(v, "Accounts", slice, err, gen.OwnAccountsCouldNotBeFound)
}

func FillOwnOrganisations[T any](v *T, acc *database.Account) {
	ownOrgs := database.OrganisationList{}
	var err error
	if acc.Role == database.HeadAdmin {
		err = ownOrgs.GetAllVisable()
	} else {
		err = ownOrgs.GetAllPartOf(acc.ID)
	}
	slice := reflect.ValueOf(ownOrgs)
	updateField(v, "Organisations", slice, err, gen.OrgNamesQueryError)
}

func updateField[T any](v *T, name string, slice reflect.Value, err error, errorText string) {
	ref := reflect.ValueOf(v).Elem()
	ref.FieldByName(name).Set(slice)
	if err != nil {
		mesg := ref.FieldByName("Message").String()
		ref.FieldByName("Message").SetString(errorText + "\n" + mesg)
	}
}
