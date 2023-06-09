package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"time"
)

var VoteDoesNotExists = "Die Abstimmung existiert nicht"

type SingleVoteViewStruct struct {
	Vote            database.Votes
	SelectedAccount string
	Accounts        database.AccountList
	Options         map[string]int
	Message         string
}

type MultipleOptionsVoteStruct struct {
	SingleVoteViewStruct
}

type RankingVoteStruct struct {
	SingleVoteViewStruct
}

type ForAgainstVoteStruct struct {
	SingleVoteViewStruct
}

func emptyVote(s *SingleVoteViewStruct, c *gin.Context, acc *database.Account) {
	if s.Vote.Info.VoteMethod == database.SingleVote {
		s.Options = map[string]int{s.Vote.Info.Options[0]: 1}
	}
	handleVote(s, c, acc)
}

func handleVote(s *SingleVoteViewStruct, c *gin.Context, acc *database.Account) {
	htmlHandler.FillOwnAccounts(s, acc)
	switch s.Vote.Info.VoteMethod {
	case database.SingleVote:
		htmlHandler.MakeSite(s, c, acc)
	case database.MultipleVotes:
		htmlHandler.MakeSite(&MultipleOptionsVoteStruct{*s}, c, acc)
	case database.VoteRanking:
		htmlHandler.MakeSite(&RankingVoteStruct{*s}, c, acc)
	case database.ThreeCategoryVoting:
		htmlHandler.MakeSite(&ForAgainstVoteStruct{*s}, c, acc)
	}
}

func GetVoteHandler(c *gin.Context) {
	f := func(s *SingleVoteViewStruct, c *gin.Context, acc *database.Account, doc *database.Document) {
		emptyVote(s, c, acc)
	}
	standardHandling(c, f)
}

func standardHandling(c *gin.Context, f func(s *SingleVoteViewStruct, c *gin.Context, acc *database.Account, doc *database.Document)) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	vote := database.Votes{}
	err := vote.GetByID(c.Query("uuid"))
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, VoteDoesNotExists)
		return
	}
	doc := database.Document{}
	if b {
		err = doc.GetByID(vote.Parent)
	} else {
		err = doc.GetByIDOnlyWithAccount(vote.Parent, acc.ID)
	}
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, VoteDoesNotExists)
		return
	}
	s := &SingleVoteViewStruct{Vote: vote}
	f(s, c, &acc, &doc)
}

func PostVoteHandler(c *gin.Context) {
	standardHandling(c, handleVotePostRequests)
}

func handleVotePostRequests(s *SingleVoteViewStruct, c *gin.Context, acc *database.Account, doc *database.Document) {
	var err error
	if doc.Type == database.RunningVote && doc.Info.Finishing.Before(time.Now().UTC()) {
		err = dataLogic.CloseDiscussionOrVote(doc.UUID)
	}
	s.SelectedAccount = generics.GetText(c, "selectedAccount")
	if err != nil {
		s.Message = ErrorWhileVoting + "\n" + s.Message
		emptyVote(s, c, acc)
		return
	}
	if !isAllowedPosting(doc, s.SelectedAccount) {
		s.Message = AccountNotForVoteAllowed + "\n" + s.Message
		emptyVote(s, c, acc)
		return
	}
	handleVoting(s, c)
	handleVote(s, c, acc)
}

var SelectedTypeNotValid = "Der übergebene type Parameter existiert nicht"

func handleVoting(s *SingleVoteViewStruct, c *gin.Context) {
	/*isValid := false
	switch true {
	case generics.GetIfType(c, "vote"):
		isValid = true
	case generics.GetIfType(c, "invalid"):
	default:
		s.Message = SelectedTypeNotValid + "\n" + s.Message
		return
	}
	var err error
	if !isValid {
		err = dataLogic.AddResultForUser(s.Vote.UUID, map[string]int{}, true, s.SelectedAccount)
	} else if !s.readIntoMap(c) {
		err = dataLogic.AddResultForUser(s.Vote.UUID, s.Options, false, s.SelectedAccount)
	}
	if err != nil {
		s.Message = ErrorWhileVoting + "\n" + s.Message
	}*/
}

var ErrorWhileVoting = "Es ist ein Fehler bei der Stimmabgabe entstanden"
var AccountNotForVoteAllowed = "Dem Account ist nicht erlaubt abzustimmen"
var SelectedOptionNotValid = "Die ausgewählten Optionen existieren nicht"
var CantSelectSameNumberTwice = "Du kannst die selbe Rangstufe nicht an zwei oder mehr Optionen vergeben"

func (s *SingleVoteViewStruct) readIntoMap(c *gin.Context) bool {
	s.Options = map[string]int{}
	switch s.Vote.Info.VoteMethod {
	case database.SingleVote:
		return s.readSingleVote(c)
	case database.MultipleVotes:
		return s.readMultipleVote(c)
	case database.VoteRanking:
		return s.readRankingVote(c)
	case database.ThreeCategoryVoting:
		return s.readThreeCategoryVote(c)
	}
	return false
}

func (s *SingleVoteViewStruct) readSingleVote(c *gin.Context) bool {
	op := generics.GetText(c, "option")
	if help.GetPositionOfString(s.Vote.Info.Options, op) == -1 {
		s.Message = SelectedOptionNotValid + "\n" + s.Message
		s.Options[s.Vote.Info.Options[0]] = 1
		return true
	}
	s.Options[op] = 1
	return false
}

func (s *SingleVoteViewStruct) readMultipleVote(c *gin.Context) bool {
	valid := false
	for _, opt := range s.Vote.Info.Options {
		if generics.GetBool(c, opt) {
			s.Options[opt] = 1
			valid = true
		}
	}
	if !valid {
		s.Message = SelectedOptionNotValid + "\n" + s.Message
		return true
	}
	return false
}

func (s *SingleVoteViewStruct) readRankingVote(c *gin.Context) bool {
	valid := false
	sameNumber := false
	selected := []int{}
	for _, opt := range s.Vote.Info.Options {
		i := generics.GetNumber(c, opt, 0, 0, s.Vote.Info.MaxPosition)
		if alreadyUsedNumber(selected, i) {
			sameNumber = true
		}
		if i != 0 {
			s.Options[opt] = i
			selected = append(selected, i)
			valid = true
		}
	}
	if !valid {
		s.Message = SelectedOptionNotValid + "\n" + s.Message
		return true
	}
	if sameNumber {
		s.Message = CantSelectSameNumberTwice + "\n" + s.Message
		return true
	}
	return false
}

func (s *SingleVoteViewStruct) readThreeCategoryVote(c *gin.Context) bool {
	valid := false
	for _, opt := range s.Vote.Info.Options {
		op := generics.GetText(c, opt)
		if op == "for" {
			s.Options[opt] = 1
			valid = true
		} else if op == "against" {
			s.Options[opt] = -1
			valid = true
		}
	}
	if !valid {
		s.Message = SelectedOptionNotValid + "\n" + s.Message
		return true
	}
	return false
}

func alreadyUsedNumber(arr []int, num int) bool {
	for _, val := range arr {
		if val == num {
			return true
		}
	}
	return false
}

func isAllowedPosting(doc *database.Document, displayName string) bool {
	/*if doc.Info.AnyPosterAllowed {
		return true
	}
	if !doc.Info.AnyPosterAllowed && !doc.Info.OrganisationPosterAllowed {
		if help.GetPositionOfString(doc.Info.Poster, displayName) != -1 {
			return true
		}
		return false
	}
	org := database.Organisation{}
	err := org.GetByName(doc.Organisation)
	if err != nil {
		return false
	}
	if help.GetPositionOfString(org.Info.Admins, displayName) != -1 ||
		help.GetPositionOfString(org.Info.User, displayName) != -1 ||
		help.GetPositionOfString(doc.Info.Poster, displayName) != -1 {
		return true
	}*/
	return false
}
