package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"zzblog/pkg/setting"
)

var db *gorm.DB

func Setup() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True",
		user, password, host, dbName))
	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.AutoMigrate(&Tag{}, &Post{})
}

func CloseDB() {
	db.Close()
}

type Page struct {
	gorm.Model
	Body        string // body
	View        int    //view count
	IsPublished string // published or not
}

// table posts
type Post struct {
	gorm.Model
	Title string //title
	Body  string //body
	View  int    //view count
	Tags  []Tag  `gorm:"many2many:post_tags;"`
}

// table tags
type Tag struct {
	gorm.Model
	Name  string // post id
	Posts []Post `gorm:"many2many:post_tags;"`
	Total int    `gorm:"-"`
}

func (page *Page) Insert() error {
	return db.Create(page).Error
}

func (page *Page) Update() error {
	return db.Model(page).Update(Page{Body: page.Body, IsPublished: page.IsPublished}).Error
}

func (page *Page) Delete() error {
	return db.Delete(page).Error
}

func GetPageById(id string) (*Page, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var page Page
	err = db.First(&page, "id=?", pid).Error
	return &page, err
}

func ListPage() ([]*Page, error) {
	var pages []*Page
	err := db.First(pages).Error
	return pages, err
}

func (post *Post) Insert() error {
	return db.Create(post).Error
}

func (post *Post) Update() error {
	p := Post{Title: post.Title}
	p.Body = post.Body
	return db.Model(post).Update(p).Error
}

func (post *Post) Delete() error {
	return db.Delete(post).Error
}

func ListPostByTagId(id string) ([]*Post, error) {
	var posts []*Post
	if len(id) == 0 {
		err := db.First(&posts).Error
		return posts, err
	} else {
		tid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, err
		}
		var tag Tag
		err = db.First(&tag, "id=?", tid).Error
		if err != nil {
			return nil, err
		}
		err = db.Model(&tag).Related(&posts, "Posts").Error
		return posts, err
	}
}

func GetPostById(id string) (*Post, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var post Post
	err = db.Preload("Tags").First(&post, "id=?", pid).Error
	return &post, err
}

func (tag *Tag) Insert() error {
	return db.FirstOrCreate(tag, "name=?", tag.Name).Error
}

func ListTag() ([]*Tag, error) {
	var tags []*Tag
	err := db.Find(&tags).Error
	return tags, err
}

func ListTagByPostId(id string) ([]*Tag, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var post Post
	var tags []*Tag
	err = db.First(&post, "id=?", pid).Error
	if err != nil {
		return nil, err
	}
	err = db.Model(&post).Related(&tags, "Tags").Error
	return tags, err
}
