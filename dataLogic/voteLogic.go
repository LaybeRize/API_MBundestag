package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/htmlHandler"
	"time"
)

var VoteAlreadyFinished htmlHandler.Message = "Die Abstimmung ist bereits beendet"
var SuccessfulVote htmlHandler.Message = "Deine Stimme wurde erfolgreich abgestimmt"
var ErrorWhileSavingVote htmlHandler.Message = "Beim Speichern der Stimme ist ein Fehler aufgetreten"
var YouAlreadyVoted htmlHandler.Message = "Du hast bereits abgestimmt"
var ErrorNotAVote htmlHandler.Message = "Dies ist keine Abstimmung"

func AddResultForUser(vote *database.Votes, m map[string]int, makeInvalid bool, user string, msg *htmlHandler.Message, positiv *bool) (err error) {
	documentLock.Lock()
	defer documentLock.Unlock()

	err = vote.GetByID(vote.UUID)
	if err != nil || vote.Finished {
		*msg = VoteAlreadyFinished + "\n" + *msg
		return
	}

	doc := database.Document{}
	err = doc.GetByID(vote.Parent)
	switch true {
	case err != nil:
		*msg = ErrorWhileSavingVote + "\n" + *msg
		return
	case doc.Type == database.FinishedVote || doc.Type == database.LegislativeText ||
		doc.Type == database.RunningDiscussion || doc.Type == database.FinishedDiscussion:
		*msg = ErrorNotAVote + "\n" + *msg
		return
	case doc.Info.Finishing.UTC().Before(time.Now().UTC()):
		vote.Finished = true //needed for the caller
		go CloseDiscussionOrVote(vote.Parent)
		*msg = VoteAlreadyFinished + "\n" + *msg
		return
	}

	_, ok := vote.Info.Results[user]
	if ok {
		*msg = YouAlreadyVoted + "\n" + *msg
		return
	}

	vote.Info.Results[user] = database.Results{
		Votee:       user,
		InvalidVote: makeInvalid,
		Votes:       m,
	}

	updatePoll(vote, user)
	err = vote.SaveChanges()
	if err != nil {
		*msg = ErrorWhileSavingVote + "\n" + *msg
		return
	}
	*positiv = true
	*msg = SuccessfulVote + "\n" + *msg
	return
}

var ForVote = "Dafür: "
var AgainstVote = "Dagegen: "
var InvalidVoteString = "Ungültige Stimme"

func updatePoll(poll *database.Votes, user string) {
	if poll.Info.Results[user].InvalidVote {
		poll.Info.Summary.InvalidVotes = append(poll.Info.Summary.InvalidVotes, user)
		poll.Info.Summary.Person[user] = InvalidVoteString
		return
	}
	switch poll.Info.VoteMethod {
	case database.SingleVote:
		vote := ""
		for key, value := range poll.Info.Results[user].Votes {
			if value != 0 {
				vote = key
			}
			poll.Info.Summary.Sums[key] += value
		}
		poll.Info.Summary.Person[user] = vote
	case database.VoteRanking:
		poll.Info.Summary.RankedMap[user] = poll.Info.Results[user].Votes
	case database.ThreeCategoryVoting:
		forVote := ""
		againstVote := ""
		for _, key := range poll.Info.Options {
			value := poll.Info.Results[user].Votes[key]
			if value == 1 {
				forVote = forVote + ", " + key
			} else if value == -1 {
				againstVote = againstVote + ", " + key
			}
			poll.Info.Summary.Sums[key] += value
		}
		if forVote != "" {
			forVote = string(([]rune(forVote))[2:])
		}
		if againstVote != "" {
			againstVote = string(([]rune(againstVote))[2:])
		}
		poll.Info.Summary.Person[user] = ForVote + forVote + "\n" + AgainstVote + againstVote
	case database.MultipleVotes:
		vote := ""
		for _, key := range poll.Info.Options {
			value := poll.Info.Results[user].Votes[key]
			if value != 0 {
				vote = vote + ", " + key
			}
			poll.Info.Summary.Sums[key] += value
		}
		if vote != "" {
			vote = string(([]rune(vote))[2:])
		}
		poll.Info.Summary.Person[user] = vote
	}
}
