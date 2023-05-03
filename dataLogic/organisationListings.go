package dataLogic

import (
	"API_MBundestag/database"
	"sort"
)

func GetNamesForSubAndMainGroups() (main []string, sub []string, err error) {
	listOrgs := database.OrganisationList{}
	err = listOrgs.GetAllMainGroups()
	if err != nil {
		return
	}
	for _, e := range listOrgs {
		main = append(main, e.MainGroup)
	}
	err = listOrgs.GetAllSubGroups()
	if err != nil {
		return
	}
	for _, e := range listOrgs {
		sub = append(sub, e.SubGroup)
	}
	return
}

func GetAllOrganisationNames() (orgs []string, err error) {
	orgList := database.OrganisationList{}
	err = orgList.GetAllVisable()
	if err != nil {
		return
	}
	for _, i := range orgList {
		orgs = append(orgs, i.Name)
	}
	err = orgList.GetAllInvisable()
	if err != nil {
		return
	}
	for _, i := range orgList {
		orgs = append(orgs, i.Name)
	}
	sort.Strings(orgs)
	return
}
