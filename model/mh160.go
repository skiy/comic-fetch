package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mh160 struct {
	Db    *gorm.DB
	Table Table
}

func (t *Mh160) CreateBook(books tb_books) (book tb_books) {
	t.Db.FirstOrCreate(&book, books)
	return book
}

func (t *Mh160) CreateChapter(chapters tb_chapter) (chapter tb_chapter) {
	t.Db.FirstOrCreate(&chapter, chapters)
	return chapter
}

func (t *Mh160) CreateImages(images tb_images) (image tb_images) {
	t.Db.FirstOrCreate(&image, images)
	return image
}
