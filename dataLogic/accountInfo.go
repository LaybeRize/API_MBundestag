package dataLogic

import "API_MBundestag/database"

func GetOrganisationList(id int64) (err error, result string) {
	list := database.OrganisationList{}
	err = list.GetAllPartOf(id)
	result = ""

	if err != nil {
		return
	}
	if list.Len() == 0 {
		return
	}

	result = list[0].Name
	for i := 1; i < list.Len(); i++ {
		result += ", " + list[i].Name
	}
	return
}

func GetTitelList(id int64) (err error, result string) {
	list := database.TitleList{}
	err = list.GetAllForUserID(id)
	result = ""

	if err != nil {
		return
	}
	if list.Len() == 0 {
		return
	}

	result = list[0].Name
	for i := 1; i < list.Len(); i++ {
		result += ", " + list[i].Name
	}
	return
}
