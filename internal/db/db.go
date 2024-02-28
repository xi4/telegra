package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Word struct {
	gorm.Model
	World       string
	Translation string
	Pages       []Page
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

func (d *DB) AddWord(word string, Translation string) uint {
	// Проверяем есть ли уже такое слово в базе
	var w Word
	d.db.Where("world = ?", word).First(&w)
	if w.ID == 0 {
		w = Word{World: word, Translation: Translation}
		d.db.Create(&w)
	}
	return w.ID
}

func (d *DB) AddPage(page Page, word_id uint) {
	page.WordID = word_id
	d.db.Create(&page)
}
