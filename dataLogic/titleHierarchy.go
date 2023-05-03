package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help"
	"sort"
	"sync"
)

type (
	TitleMainGroupArray []TitleMainGroup
	TitleMainGroup      struct {
		Name   string
		Groups []TitleSubGroup
	}
	TitleSubGroup struct {
		Name   string
		Titles []database.Title
	}
)

var titleHierarchyLock = sync.Mutex{}
var titleHierarchy TitleMainGroupArray
var mainGroupNames []string
var subGroupNames []string
var titleNames []string

func GetTitleHierarchy() TitleMainGroupArray {
	titleHierarchyLock.Lock()
	defer titleHierarchyLock.Unlock()
	return titleHierarchy
}

func GetMainGroupNames() []string {
	titleHierarchyLock.Lock()
	defer titleHierarchyLock.Unlock()
	return mainGroupNames
}

func GetSubGroupNames() []string {
	titleHierarchyLock.Lock()
	defer titleHierarchyLock.Unlock()
	return subGroupNames
}

func GetTitleNames() []string {
	titleHierarchyLock.Lock()
	defer titleHierarchyLock.Unlock()
	return titleNames
}

func RefreshTitleHierarchy() (err error) {
	titleHierarchyLock.Lock()
	defer titleHierarchyLock.Unlock()
	//Empty arrays
	mainGroupNames = []string{}
	subGroupNames = []string{}
	titleNames = []string{}

	//First fetch all Organisations
	list := database.TitleList{}
	err = list.GetAll()
	//If list is empty just create an empty tempTitleHierarchy
	if err != nil {
		return
	}
	if len(list) == 0 {
		titleHierarchy = []TitleMainGroup{}
		return nil
	}
	//If list is only one long, just fill in the correct parameter and return
	main := database.MainGroupListTitle(list)

	if main.Len() == 1 {

		mainGroupNames = []string{main[0].MainGroup}
		subGroupNames = []string{main[0].SubGroup}
		titleNames = []string{main[0].Name}

		titleHierarchy = []TitleMainGroup{
			{
				Name: main[0].MainGroup,
				Groups: []TitleSubGroup{{
					Name:   main[0].SubGroup,
					Titles: main,
				}},
			},
		}
		return
	}
	//Otherwise sort the array, so that equal names are always after each other

	sort.Sort(main)

	lastPos := 0
	var mainGroupList []database.SubGroupListTitle
	//Then split the array into its blocks
	for i := 0; i < main.Len()-1; i++ {
		titleNames = append(titleNames, main[i].Name)
		if main[i].MainGroup != main[i+1].MainGroup {
			mainGroupList = append(mainGroupList, database.SubGroupListTitle(main[lastPos:i+1]))
			lastPos = i + 1
		}
	}
	titleNames = append(titleNames, main[main.Len()-1].Name)

	//add the last block
	mainGroupList = append(mainGroupList, database.SubGroupListTitle(main[lastPos:]))
	var tempTitleHierarchy []TitleMainGroup
	//Fill the tempTitleHierarchy
	for _, array := range mainGroupList {
		//Either with the correct Element directly
		mainGroupNames = append(mainGroupNames, array[0].MainGroup)

		if array.Len() == 1 {
			subGroupNames = append(subGroupNames, array[0].SubGroup)

			tempTitleHierarchy = append(tempTitleHierarchy, TitleMainGroup{
				Name: array[0].MainGroup,
				Groups: []TitleSubGroup{{
					Name:   array[0].SubGroup,
					Titles: array,
				},
				},
			})
			continue
		}
		//or of it is not clear, that only one subgroup exists, just fill it into a blank first
		sort.Sort(array)
		tempTitleHierarchy = append(tempTitleHierarchy, TitleMainGroup{
			Name: array[0].MainGroup,
			Groups: []TitleSubGroup{{
				Name:   "",
				Titles: array,
			},
			},
		})
	}
	//Clear actual hierarchy
	titleHierarchy = TitleMainGroupArray{}
	//Find the not assigned subgroups and repeat the process from the main groups one step lower
	for _, mainGroup := range tempTitleHierarchy {
		if mainGroup.Groups[0].Name != "" {
			//if the main group is already well-defined add it to the actual TitleHierarchy
			titleHierarchy = append(titleHierarchy, mainGroup)
			continue
		}
		//otherwise split the titles into it's subgroups
		lastPos = 0
		titles := database.TitleList(mainGroup.Groups[0].Titles)
		var titleListList []database.TitleList
		//Then split the array into its blocks
		for i := 0; i < titles.Len()-1; i++ {
			if titles[i].SubGroup != titles[i+1].SubGroup {
				titleListList = append(titleListList, titles[lastPos:i+1])
				lastPos = i + 1
			}
		}
		titleListList = append(titleListList, titles[lastPos:])
		//Then define main group well
		for num, array := range titleListList {
			subGroupNames = append(subGroupNames, array[0].SubGroup)

			sort.Sort(&array)
			if num == 0 {
				mainGroup.Groups[0].Name = array[0].SubGroup
				mainGroup.Groups[0].Titles = array
				continue
			}
			mainGroup.Groups = append(mainGroup.Groups, TitleSubGroup{
				Name:   array[0].SubGroup,
				Titles: array,
			})
		}
		//and add it to the actual TitleHierarchy
		titleHierarchy = append(titleHierarchy, mainGroup)
	}
	subGroupNames = help.RemoveDuplicates(subGroupNames)
	return
}
