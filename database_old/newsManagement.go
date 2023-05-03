package database

import (
	"database/sql"
	"log"
	"time"
)

var NewsSchema = `
CREATE TABLE IF NOT EXISTS publications (
    uuid text UNIQUE NOT NULL,
    creation_time TIMESTAMP NOT NULL,
    publication_time TIMESTAMP NOT NULL,
    publicated BOOLEAN NOT NULL,
    hast BOOLEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS article (
    uuid text UNIQUE NOT NULL,
    publication text NOT NULL,
    written TIMESTAMP NOT NULL,
    author TEXT NOT NULL,
    flair TEXT NOT NULL,
    headline TEXT NOT NULL,
    subtitle TEXT,
    content TEXT NOT NULL,
    html_content TEXT NOT NULL
);
`

func TestNewsDB() {
	TestDatabase("DROP TABLE IF EXISTS publications, article;", "")
	InitNewsDatabase()

	pub := Publication{
		UUID:         EternatityPublicationName,
		PublishTime:  time.Now().UTC(),
		Publicated:   false,
		BreakingNews: false,
	}

	err := pub.CreateMe()
	if err != nil {
		log.Fatal(err)
	}
}

var EternatityPublicationName = "theFirstOfThemAll"

type (
	PublicationList []Publication
	Publication     struct {
		UUID         string
		CreateTime   time.Time `db:"creation_time"`
		PublishTime  time.Time `db:"publication_time"`
		Publicated   bool
		BreakingNews bool `db:"hast"`
	}
	ArticleList []Article
	Article     struct {
		UUID        string
		Publication string
		Written     time.Time
		Author      string
		Flair       string
		Headline    string
		Subtitle    sql.NullString
		Content     string
		HTMLContent string `db:"html_content"`
	}
)

func InitNewsDatabase() {
	DB.MustExec(NewsSchema)
}

func (publication *Publication) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO publications (uuid, creation_time, publication_time, publicated, hast) VALUES (:uuid, :creation_time, :publication_time, :publicated, :hast)", map[string]interface{}{
		"uuid":             publication.UUID,
		"creation_time":    time.Now().UTC(),
		"publication_time": publication.PublishTime,
		"publicated":       publication.Publicated,
		"hast":             publication.BreakingNews,
	})
	return
}

func (publication *Publication) GetByID(uuid string) (err error) {
	err = DB.Get(publication, "SELECT * FROM publications WHERE uuid=$1;", uuid)
	return
}

func (publication *Publication) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE publications SET creation_time=:creation_time, publication_time=:publication_time, publicated=:publicated, hast=:hast WHERE uuid=:uuid", map[string]interface{}{
		"uuid":             publication.UUID,
		"creation_time":    time.Now().UTC(),
		"publication_time": publication.PublishTime,
		"publicated":       publication.Publicated,
		"hast":             publication.BreakingNews,
	})
	return
}

func (publication *Publication) DeleteMe() (err error) {
	_, err = DB.NamedExec("DELETE FROM publications WHERE uuid=:uuid", map[string]interface{}{
		"uuid": publication.UUID,
	})
	return
}

func (publicationList *PublicationList) GetOnlyUnpublicated() (err error) {
	err = DB.Select(publicationList, "SELECT * FROM publications WHERE publicated = false ORDER BY creation_time;")
	return
}

func (publicationList *PublicationList) GetPublicationBeforeDate(date time.Time, amount int) (err error) {
	// Comment: If you want the day itself, you have to set the time to 23:59:59
	err = DB.Select(publicationList, "SELECT * FROM publications WHERE publication_time < $1 AND publicated = true ORDER BY publication_time DESC LIMIT $2;", date, amount)
	return
}

func (publicationList *PublicationList) GetPublicationAfter(publicationUUID string, amount int) (err error, exists bool) {
	exists = true
	pub := Publication{}
	err = pub.GetByID(publicationUUID)
	if err == sql.ErrNoRows {
		exists = false
		pub.PublishTime = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = DB.Select(publicationList, "SELECT * FROM publications WHERE publication_time < $2 AND uuid != $1 AND publicated = true ORDER BY publication_time DESC LIMIT $3;", publicationUUID, pub.PublishTime, amount)
	return
}

func (publicationList *PublicationList) GetPublicationBefore(publicationUUID string, amount int) (err error, exists bool) {
	exists = true
	pub := Publication{}
	err = pub.GetByID(publicationUUID)
	if err == sql.ErrNoRows {
		exists = false
		pub.PublishTime = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = DB.Select(publicationList, "SELECT * FROM (SELECT * FROM publications WHERE publication_time > $2 AND uuid != $1 AND publicated = true ORDER BY publication_time LIMIT $3) as X ORDER BY X.publication_time DESC;", publicationUUID, pub.PublishTime, amount)
	return
}

func (article *Article) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO article (uuid, publication, written, author, flair, headline, subtitle, content, html_content) VALUES (:uuid, :publication, :written, :author, :flair, :headline, :subtitle, :content, :html_content)", map[string]interface{}{
		"uuid":         article.UUID,
		"publication":  article.Publication,
		"written":      time.Now().UTC(),
		"author":       article.Author,
		"flair":        article.Flair,
		"headline":     article.Headline,
		"subtitle":     article.Subtitle,
		"content":      article.Content,
		"html_content": article.HTMLContent,
	})
	return
}

func (article *Article) GetByID(uuid string) (err error) {
	err = DB.Get(article, "SELECT * FROM article WHERE uuid=$1;", uuid)
	return
}

func (article *Article) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE article SET publication=:publication, author=:author, flair=:flair, headline=:headline, subtitle=:subtitle, content=:content, html_content=:html_content WHERE uuid=:uuid", map[string]interface{}{
		"uuid":         article.UUID,
		"publication":  article.Publication,
		"author":       article.Author,
		"flair":        article.Flair,
		"headline":     article.Headline,
		"subtitle":     article.Subtitle,
		"content":      article.Content,
		"html_content": article.HTMLContent,
	})
	return
}

func (article *Article) DeleteMe() (err error) {
	_, err = DB.NamedExec("DELETE FROM article WHERE uuid=:uuid", map[string]interface{}{
		"uuid": article.UUID,
	})
	return
}

func (articleList *ArticleList) GetAllArticlesToPublication(publicationUUID string) (err error) {
	err = DB.Select(articleList, "SELECT * FROM article WHERE publication=$1 ORDER BY written;", publicationUUID)
	return
}
