package dataLogic

import "API_MBundestag/database"

func CloseDiscussionOrVote(uuid string) (err error) {
	documentLock.Lock()
	defer documentLock.Unlock()
	temp := database.Document{}
	err = temp.GetByID(uuid)
	if err != nil {
		return
	}
	switch temp.Type {
	case database.RunningDiscussion:
		temp.Type = database.FinishedDiscussion
	case database.RunningVote:
		for _, str := range temp.Info.Votes {
			err = closePoll(str)
		}
		temp.Type = database.FinishedVote
	}
	err = temp.SaveChanges()
	return
}

func closePoll(uuid string) (err error) {
	vote := database.Votes{}
	err = vote.GetByID(uuid)
	if err != nil || vote.Finished {
		return
	}

	vote.Finished = true
	createCSV(&vote)
	err = vote.SaveChanges()
	return
}

func createCSV(poll *database.Votes) {

}
