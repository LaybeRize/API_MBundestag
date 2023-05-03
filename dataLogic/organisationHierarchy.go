package dataLogic

import (
	"API_MBundestag/database"
	"sort"
)

type (
	OrganisationMainGroupArray []OrganisationMainGroup
	OrganisationMainGroup      struct {
		Name   string
		Amount int
		Groups []OrganisationSubGroup
	}
	OrganisationSubGroup struct {
		Name          string
		Amount        int
		Organisations []database.Organisation
	}
)

func (OrgHierarchy *OrganisationMainGroupArray) GetOrganisationHierarchy(acc database.Account, hiddenMode bool) (err error) {
	//First fetch all Organisations
	list := database.OrganisationList{}
	if hiddenMode && (acc.Role == database.Admin || acc.Role == database.HeadAdmin) {
		err = list.GetAllInvisable()
	} else if acc.Role == database.Admin || acc.Role == database.HeadAdmin {
		err = list.GetAllVisable()
	} else {
		err = list.GetAllVisibleFor(acc.ID)
	}

	//If list is empty just create an empty tempOrgHierarchy
	if err != nil {
		return
	}
	if len(list) == 0 {
		*OrgHierarchy = []OrganisationMainGroup{}
		return nil
	}
	//If list is only one long, just fill in the correct parameter and return
	main := database.MainGroupListOrg(list)
	if main.Len() == 1 {
		*OrgHierarchy = []OrganisationMainGroup{
			{
				Name: main[0].MainGroup,
				Groups: []OrganisationSubGroup{{
					Name:          main[0].SubGroup,
					Amount:        1,
					Organisations: main,
				}},
			},
		}
		return
	}
	//Otherwise sort the array, so that equal names are always after each other
	sort.Sort(main)
	lastPos := 0
	var subList []database.SubGroupListOrg
	//Then split the array into its blocks
	for i := 0; i < main.Len()-1; i++ {
		if main[i].MainGroup != main[i+1].MainGroup {
			subList = append(subList, database.SubGroupListOrg(main[lastPos:i+1]))
			lastPos = i + 1
		}
	}
	//add the last block
	subList = append(subList, database.SubGroupListOrg(main[lastPos:]))
	var tempOrgHierarchy []OrganisationMainGroup
	//Fill the tempOrgHierarchy
	for _, array := range subList {
		//Either with the correct Element directly
		if array.Len() == 1 {
			tempOrgHierarchy = append(tempOrgHierarchy, OrganisationMainGroup{
				Name: array[0].MainGroup,
				Groups: []OrganisationSubGroup{{
					Name:          array[0].SubGroup,
					Amount:        1,
					Organisations: array,
				},
				},
			})
			continue
		}
		//or of it is not clear, that only one subgroup exists, just fill it into a blank first
		sort.Sort(array)
		tempOrgHierarchy = append(tempOrgHierarchy, OrganisationMainGroup{
			Name: array[0].MainGroup,
			Groups: []OrganisationSubGroup{{
				Name:          "",
				Amount:        len(array),
				Organisations: array,
			},
			},
		})
	}
	//Clear actual hierarchy
	*OrgHierarchy = OrganisationMainGroupArray{}
	//Find the not assigned subgroups and repeat the process from the main groups one step lower
	for _, mainGroup := range tempOrgHierarchy {
		if mainGroup.Groups[0].Name != "" {
			//if the main group is already well-defined add it to the actual OrgHierarchy
			*OrgHierarchy = append(*OrgHierarchy, mainGroup)
			continue
		}
		//otherwise split the orgs into it's subgroups
		lastPos = 0
		orgs := database.OrganisationList(mainGroup.Groups[0].Organisations)
		var orgListList []database.OrganisationList
		//Then split the array into its blocks
		for i := 0; i < orgs.Len()-1; i++ {
			if orgs[i].SubGroup != orgs[i+1].SubGroup {
				orgListList = append(orgListList, orgs[lastPos:i+1])
				lastPos = i + 1
			}
		}
		orgListList = append(orgListList, orgs[lastPos:])
		//Then do the well-defining of the maingroup
		for num, array := range orgListList {
			sort.Sort(&array)
			if num == 0 {
				mainGroup.Groups[0].Name = array[0].SubGroup
				mainGroup.Groups[0].Organisations = array
				mainGroup.Groups[0].Amount = array.Len()
				continue
			}
			mainGroup.Groups = append(mainGroup.Groups, OrganisationSubGroup{
				Name:          array[0].SubGroup,
				Amount:        len(array),
				Organisations: array,
			})
		}
		//and add it to the actual OrgHierarchy
		*OrgHierarchy = append(*OrgHierarchy, mainGroup)
	}
	return
}

func (OrgHierarchy *OrganisationMainGroupArray) SetAmountForMainGroup() {
	for i, element := range *OrgHierarchy {
		sum := 0
		for _, sub := range element.Groups {
			sum += sub.Amount
		}
		(*OrgHierarchy)[i].Amount = sum
	}
}
