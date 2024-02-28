package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Word struct {
	gorm.Model
	World string
	Pages []Page
}

type Page struct {
	gorm.Model
	Title   string
	Address string
	Images  []Image
	Videos  []Video
	WordID  uint
}

type Image struct {
	gorm.Model
	Url    string
	PageID uint
}

type Video struct {
	gorm.Model
	Url    string
	PageID uint
}

type DB struct {
	db *gorm.DB
}

func New() *DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Word{}, &Page{}, &Image{}, &Video{})
	return &DB{db: db}
}

func (d *DB) AddWord(word string, pages []Page) {
	d.db.Create(&Word{World: word, Pages: pages})
}

func (d *DB) AddPage(page Page) {
	d.db.Create(&page)
}

func (d *DB) AddImage(image Image) {
	d.db.Create(&image)
}

func (d *DB) AddVideo(video Video) {
	d.db.Create(&video)
}

func (d *DB) GetWord(word string) Word {
	var w Word
	d.db.Where("world = ?", word).First(&w)
	return w
}

func (d *DB) GetPage(id uint) Page {
	var p Page
	d.db.First(&p, id)
	return p
}

func (d *DB) GetImage(id uint) Image {
	var i Image
	d.db.First(&i, id)
	return i
}

func (d *DB) GetVideo(id uint) Video {
	var v Video
	d.db.First(&v, id)
	return v
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) GetPages() []Page {
	var pages []Page
	d.db.Find(&pages)
	return pages
}

func (d *DB) GetImages() []Image {
	var images []Image
	d.db.Find(&images)
	return images
}

func (d *DB) GetVideos() []Video {
	var videos []Video
	d.db.Find(&videos)
	return videos
}

func (d *DB) GetWords() []Word {
	var words []Word
	d.db.Find(&words)
	return words
}

func (d *DB) GetPagesByWord(word string) []Page {
	var pages []Page
	d.db.Where("word_id = ?", word).Find(&pages)
	return pages
}

func (d *DB) GetImagesByPage(page Page) []Image {
	var images []Image
	d.db.Where("page_id = ?", page.ID).Find(&images)
	return images
}

func (d *DB) GetVideosByPage(page Page) []Video {
	var videos []Video
	d.db.Where("page_id = ?", page.ID).Find(&videos)
	return videos
}

func (d *DB) GetWordByPage(page Page) Word {
	var word Word
	d.db.First(&word, page.WordID)
	return word
}

func (d *DB) GetPageByImage(image Image) Page {
	var page Page
	d.db.First(&page, image.PageID)
	return page
}

func (d *DB) GetPageByVideo(video Video) Page {
	var page Page
	d.db.First(&page, video.PageID)
	return page
}

func (d *DB) GetWordByImage(image Image) Word {
	var page Page
	d.db.First(&page, image.PageID)
	return d.GetWordByPage(page)
}

func (d *DB) GetWordByVideo(video Video) Word {
	var page Page
	d.db.First(&page, video.PageID)
	return d.GetWordByPage(page)
}

func (d *DB) GetImagesByWord(word string) []Image {
	var images []Image
	pages := d.GetPagesByWord(word)
	for _, page := range pages {
		images = append(images, d.GetImagesByPage(page)...)
	}
	return images
}

func (d *DB) GetVideosByWord(word string) []Video {
	var videos []Video
	pages := d.GetPagesByWord(word)
	for _, page := range pages {
		videos = append(videos, d.GetVideosByPage(page)...)
	}
	return videos
}

func (d *DB) GetPagesByImage(image Image) []Page {
	var pages []Page
	d.db.Where("id = ?", image.PageID).Find(&pages)
	return pages
}

func (d *DB) GetPagesByVideo(video Video) []Page {
	var pages []Page
	d.db.Where("id = ?", video.PageID).Find(&pages)
	return pages
}
