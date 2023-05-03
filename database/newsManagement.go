package database

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

var EternatityPublicationName = "theFirstOfThemAll"

type (
	PublicationList []Publication
	Publication     struct {
		UUID         string    `gorm:"primaryKey"`
		CreateTime   time.Time `gorm:"column:creation_time"`
		PublishTime  time.Time `gorm:"column:publication_time"`
		Publicated   bool
		BreakingNews bool `gorm:"column:hast"`
	}
	ArticleList []Article
	Article     struct {
		UUID        string `gorm:"primaryKey"`
		Publication string
		Written     time.Time
		Author      string
		Flair       string
		Headline    string
		Subtitle    sql.NullString
		Content     string
		HTMLContent string `gorm:"column:html_content"`
	}
)

func (publication *Publication) CreateMe() (err error) {
	publication.CreateTime = time.Now().UTC()
	err = db.Create(publication).Error
	return
}

func (publication *Publication) GetByID(uuid string) (err error) {
	*publication = Publication{}
	err = db.Where("uuid=?", uuid).First(publication).Error
	return
}

func (publication *Publication) SaveChanges() (err error) {
	publication.CreateTime = time.Now().UTC()
	err = db.Save(publication).Error
	return
}

func (publication *Publication) DeleteMe() (err error) {
	err = db.Delete(publication).Error
	return
}

func (publicationList *PublicationList) GetOnlyUnpublicated() (err error) {
	*publicationList = PublicationList{}
	err = db.Where("publicated = false").Order("creation_time").Find(publicationList).Error
	return
}

func (publicationList *PublicationList) GetPublicationBeforeDate(date time.Time, amount int) (err error) {
	// Comment: If you want the day itself, you have to set the time to 23:59:59
	*publicationList = PublicationList{}
	err = db.Where(" publication_time < ? AND publicated = true", date).Order("publication_time DESC").Limit(amount).Find(publicationList).Error
	return
}

func (publicationList *PublicationList) GetPublicationAfter(publicationUUID string, amount int) (err error, exists bool) {
	*publicationList = PublicationList{}
	exists = true
	pub := Publication{}
	err = pub.GetByID(publicationUUID)
	if err == gorm.ErrRecordNotFound {
		exists = false
		pub.PublishTime = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = getBasicNewsQuery(publicationUUID, amount).Where("publication_time < ?", pub.PublishTime).Order("publication_time desc").Find(publicationList).Error
	return
}

func (publicationList *PublicationList) GetPublicationBefore(publicationUUID string, amount int) (err error, exists bool) {
	*publicationList = PublicationList{}
	exists = true
	pub := Publication{}
	err = pub.GetByID(publicationUUID)
	if err == gorm.ErrRecordNotFound {
		exists = false
		pub.PublishTime = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = db.Select("*").Table("(?) as X", getBasicNewsQuery(publicationUUID, amount).Order("publication_time").Where("publication_time > ?", pub.PublishTime)).Order("X.publication_time desc").Find(publicationList).Error
	return
}

func getBasicNewsQuery(uuid string, amount int) *gorm.DB {
	return db.Select("*").Where("uuid != ? AND publicated = true", uuid).Limit(amount).Table("publications")
}

func (article *Article) CreateMe() (err error) {
	article.Written = time.Now().UTC()
	err = db.Create(article).Error
	return
}

func (article *Article) GetByID(uuid string) (err error) {
	*article = Article{}
	err = db.Where("uuid = ?", uuid).First(article).Error
	return
}

func (article *Article) SaveChanges() (err error) {
	err = db.Save(article).Error
	return
}

func (article *Article) DeleteMe() (err error) {
	err = db.Delete(article).Error
	return
}

func (articleList *ArticleList) GetAllArticlesToPublication(publicationUUID string) (err error) {
	*articleList = ArticleList{}
	err = db.Where("publication=?", publicationUUID).Order("written").Find(articleList).Error
	return
}

func (publication *Publication) UpdateAllArticles(newUUID string) (err error) {
	err = db.Model(&Article{}).Where("publication=?", publication.UUID).Select("publication").Updates(&Article{Publication: newUUID}).Error
	return
}
