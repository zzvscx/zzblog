package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
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
	db.AutoMigrate(&Tag{}, &Page{}, &Post{}, &User{})
}

func CloseDB() {
	db.Close()
}

type Page struct {
	gorm.Model
	Title       string //title
	Body        string // body
	View        int    //view count
	IsPublished string // published or not
}

// table posts
type Post struct {
	gorm.Model
	Title       string //title
	Body        string //body
	View        int    //view count
	IsPublished string // published or not
	Tags        []*Tag `gorm:"many2many:post_tags;"`
}

// table tags
type Tag struct {
	gorm.Model
	Name  string  // post id
	Total int     `gorm:"-"`
	Posts []*Post `gorm:"many2many:post_tags;"`
}

type QrArchive struct {
	ArchiveDate time.Time
	Total       int
	Year        int
	Month       int
}

type User struct {
	gorm.Model
	Email       string `gorm:"unique_index"` //邮箱
	Password    string //密码
	VerifyState string //邮箱验证状态
	IsAdmin     bool   //是否是管理员
}

func (page *Page) Insert() error {
	return db.Create(page).Error
}

func (page *Page) Update() error {
	return db.Model(page).Update(Page{Title: page.Title, Body: page.Body, IsPublished: page.IsPublished}).Error
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
	err := db.Related("Tags").Find(pages).Error
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
		err := db.Find(&posts).Error
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
	rows, err := db.Raw("select t.id, t.name, count(*) from tag t inner join post_tags pt on t.id=pt.tag_id " +
		"inner join post p on pt.post_id=p.id group by pt.tag_id").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		rows.Scan(&tag.ID, &tag.Name, &tag.Total)
		tags = append(tags, &tag)
	}

	return tags, nil
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

func PostCountByArchives() ([]*QrArchive, error) {
	var archives []*QrArchive
	sql := `select DATE_FORMAT(created_at,'%Y-%m') as month, count(*) as total from post group by month order by month desc`
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var archive QrArchive
		var month string
		rows.Scan(&month, &archive.Total)
		archive.ArchiveDate, _ = time.Parse("2006-01", month)
		archive.Year = archive.ArchiveDate.Year()
		archive.Month = int(archive.ArchiveDate.Month())
		archives = append(archives, &archive)
	}
	return archives, nil

}

func ListPostByArchive(year, month string) ([]*Post, error) {
	if len(month) == 1 {
		month = "0" + month
	}
	condition := fmt.Sprintf("%s-%s", year, month)
	rows, err := db.Raw("select title,body from post where date_format(created_at,'%Y-%m') = ?", condition).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, 0)
	for rows.Next() {
		var post Post
		rows.Scan(&post.Title, &post.Body)
		posts = append(posts, &post)
	}
	return posts, nil
}

func (user *User) Insert() error {
	return db.Create(user).Error
}

func (user *User) Update() error {
	return db.Save(user).Error
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.First(&user, "email = ?", username).Error
	return &user, err
}
