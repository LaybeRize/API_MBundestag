package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var VoteSchema = `
CREATE TABLE IF NOT EXISTS votes (
    uuid text UNIQUE NOT NULL,
    parent TEXT NOT NULL,
    question TEXT NOT NULL,
    wv_number BOOLEAN NOT NULL,
    wv_names BOOLEAN NOT NULL,
    na_names BOOLEAN NOT NULL,
    finished BOOLEAN NOT NULL,
    info jsonb NOT NULL
);
`

func TestVotesDB() {
	TestDatabase("DROP TABLE IF EXISTS votes;", "")
	InitDocumentsDatabase()
}

type (
	VotesList []Votes
	Votes     struct {
		UUID                   string
		Parent                 string
		Question               string
		ShowNumbersWhileVoting bool `db:"wv_number"`
		ShowNamesWhileVoting   bool `db:"wv_names"`
		ShowNamesAfterVoting   bool `db:"na_names"`
		Finished               bool
		Info                   VoteInfo
	}
	VoteInfo struct {
		Allowed     []string           `json:"allowed"`
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

func InitVotesDatabase() {
	DB.MustExec(VoteSchema)
}

func (docI *VoteInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &docI)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &docI)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (docI *VoteInfo) Value() driver.Value {
	l, _ := json.Marshal(&docI)
	return l
}

func (votes *Votes) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO votes (uuid, parent, question, wv_number, wv_names, na_names, finished, info) VALUES (:uuid, :parent, :question, :wv_number, :wv_names, :na_names, :finished, :info)", map[string]interface{}{
		"uuid":      votes.UUID,
		"parent":    votes.Parent,
		"question":  votes.Question,
		"wv_number": votes.ShowNumbersWhileVoting,
		"wv_names":  votes.ShowNamesWhileVoting,
		"na_names":  votes.ShowNamesAfterVoting,
		"finished":  false,
		"info":      votes.Info.Value(),
	})
	return
}

func (votes *Votes) GetByID(uuid string) (err error) {
	err = DB.Get(votes, "SELECT * FROM votes WHERE uuid=$1;", uuid)
	return
}

func (votes *Votes) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE votes SET info=:info, finished=:finished WHERE uuid=:uuid", map[string]interface{}{
		"uuid":     votes.UUID,
		"finished": votes.Finished,
		"info":     votes.Info.Value(),
	})
	return
}

func (voteList *VotesList) GetVotesForDocument(docUUID string) (err error) {
	err = DB.Select(voteList, "SELECT FROM votes WHERE parent=$1 ORDER BY question", docUUID)
	return
}
