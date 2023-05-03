package htmlPress

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type NewspaperListViewStruct struct {
	Search             bool
	HasNext            bool
	HasBefore          bool
	NextUUID           string
	BeforeUUID         string
	Amount             int
	BreakingNewsFormat string
	NormalNewsFormat   string
	PubList            database.PublicationList
}

type NewspaperHiddenListViewStruct struct {
	NewspaperListViewStruct
}

func GetNewsPaperHiddenListPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	listStruct := NewspaperListViewStruct{}
	err := listStruct.PubList.GetOnlyUnpublicated()
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLoadingNewsPaper)
		return
	}

	listStruct.BreakingNewsFormat = generics.FormatHiddenBreakingNews
	listStruct.NormalNewsFormat = generics.FormatHiddenNormalNews
	htmlHandler.MakeSite(&NewspaperHiddenListViewStruct{listStruct}, c, &acc)
}

func GetNewsPaperListPage(c *gin.Context) {
	acc, _ := dataLogic.CheckUserPrivileged(c)

	i := htmlHandler.ExtractAmount(c, 1, 50, 20)

	var err error
	listStruct := &NewspaperListViewStruct{}
	if c.Query("type") == "before" {
		err = listStruct.validateNewsPaperReadPageBefore(c, i)
	} else {
		err = listStruct.validateNewsPaperReadNextPage(c, i)
	}
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLoadingNewsPaper)
		return
	}

	listStruct.Amount = i
	listStruct.BreakingNewsFormat = generics.FormatBreakingNews
	listStruct.NormalNewsFormat = generics.FormatNormalNews
	listStruct.Search = true
	htmlHandler.MakeSite(listStruct, c, &acc)
}

func (listStruct *NewspaperListViewStruct) validateNewsPaperReadNextPage(c *gin.Context, i int) error {
	listStruct.PubList = database.PublicationList{}
	err, exists := listStruct.PubList.GetPublicationAfter(c.Query("uuid"), i+1)
	if len(listStruct.PubList) == 0 {
		return err
	}
	if len(listStruct.PubList) == i+1 {
		listStruct.HasNext = true
		listStruct.NextUUID = listStruct.PubList[i-1].UUID
		listStruct.PubList = listStruct.PubList[:i]
	}
	if exists {
		listStruct.HasBefore = true
		listStruct.BeforeUUID = listStruct.PubList[0].UUID
	}
	return err
}

func (listStruct *NewspaperListViewStruct) validateNewsPaperReadPageBefore(c *gin.Context, i int) error {
	listStruct.PubList = database.PublicationList{}
	err, exists := listStruct.PubList.GetPublicationBefore(c.Query("uuid"), i+1)
	if len(listStruct.PubList) == 0 {
		return err
	}
	if len(listStruct.PubList) == i+1 {
		listStruct.HasBefore = true
		listStruct.PubList = listStruct.PubList[1:]
		listStruct.BeforeUUID = listStruct.PubList[0].UUID
	}
	if exists {
		listStruct.HasNext = true
		listStruct.NextUUID = listStruct.PubList[len(listStruct.PubList)-1].UUID
	}
	return err
}
