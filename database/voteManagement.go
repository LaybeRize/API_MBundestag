package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

func (voteI *VoteInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &voteI)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &voteI)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (voteI *VoteInfo) Value() driver.Value {
	l, _ := json.Marshal(&voteI)
	return l
}

type (
	VotesList []Votes
	Votes     struct {
		UUID                   string `gorm:"primaryKey"`
		Parent                 string
		Question               string
		ShowNumbersWhileVoting bool
		ShowNamesWhileVoting   bool
		ShowNamesAfterVoting   bool
		Finished               bool
		Info                   VoteInfo `gorm:"type:jsonb"`
	}
	VoteInfo struct {
		Results     map[string]Results `json:"results"`
		Summary     Summary            `json:"summary"`
		VoteMethod  VoteType           `json:"voteMethod"`
		MaxPosition int                `json:"maxPosition"`
		Options     []string           `json:"options"`
	}
	Results struct {
		Votee       string         `json:"votee"`
		InvalidVote bool           `json:"invald"`
		Votes       map[string]int `json:"votes"`
	}
	Summary struct {
		Sums         map[string]int            `json:"sums"`
		RankedMap    map[string]map[string]int `json:"rankedMap"`
		Person       map[string]string         `json:"person"` //the option the person voted for
		InvalidVotes []string                  `json:"invalidVotes"`
		CSV          string                    `json:"csv"` //saves the data as a CSV for the ranked Map
	}
)

const (
	SingleVote          VoteType = "single_vote"
	MultipleVotes       VoteType = "multiple_votes"
	VoteRanking         VoteType = "vote_ranking"
	ThreeCategoryVoting VoteType = "three_category_voting" //for against neutral
)

var VoteTypes = []VoteType{SingleVote, MultipleVotes, VoteRanking, ThreeCategoryVoting}
var VoteTranslation = map[VoteType]string{
	SingleVote:          "Einzelstimmenwahl",
	MultipleVotes:       "Mehrstimmenwahl",
	VoteRanking:         "Gewichtete Wahl",
	ThreeCategoryVoting: "Dafür-Dagegen-Enthaltung-Wahl",
}

func (votes *Votes) CreateMe() (err error) {
	err = db.Create(votes).Error
	return
}

func (votes *Votes) GetByID(uuid string) (err error) {
	*votes = Votes{}
	err = db.First(votes, "uuid = ?", uuid).Error
	return
}

func (votes *Votes) SaveChanges() (err error) {
	err = db.Save(votes).Error
	return
}
