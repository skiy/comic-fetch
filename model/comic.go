package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Comic struct {
	Db    *gorm.DB
	Table Table
}

func (t *Comic) CreateBook(books tb_books) (book tb_books) {
	t.Db.FirstOrCreate(&book, books)
	return book
}

func (t *Comic) CreateChapter(chapters tb_chapter) (chapter tb_chapter) {
	t.Db.FirstOrCreate(&chapter, chapters)
	return chapter
}

func (t *Comic) CreateImages(images tb_images) (image tb_images) {
	t.Db.FirstOrCreate(&image, images)
	return image
}

func (t *Comic) GetBookList(id int) (books []tb_books) {
	if id != 0 {
		t.Db.Where("id = ?", id)
	}
	t.Db.Find(&books)
	return
}

func (t *Comic) FetchImageList() (images []FtImages) {
	//sql := "SELECT i.bid,i.cid,i.order_id,i.origin_url,b.name,b.origin_url,b.origin_flag FROM tb_images AS i LEFT JOIN tb_books b ON b.id = i.bid WHERE I.image_url = '' ORDER BY i.id ASC" .Limit(10)
	t.Db.Table("tb_images AS i").Select("i.id,i.bid,i.cid,i.order_id,i.origin_url AS image_url,b.name,b.origin_url,b.origin_flag").Joins("LEFT JOIN tb_books b ON b.id = i.bid").Where("i.image_url = ''").Order("i.id ASC").Scan(&images)
	return
}

func (t *Comic) GetChapterList(bid int) (chapters []tb_chapter) {
	if bid != 0 {
		t.Db.Where("bid = ?", bid)
	}
	t.Db.Find(&chapters)
	return
}

/**
更新漫画
*/
func (t *Comic) UpdateBook(id int, book tb_books) bool {
	books := new(tb_books)
	t.Db.Model(&books).Where("id = ?", id).UpdateColumns(book)

	return true
}

/**
更新图片
*/
func (t *Comic) UpdateImage(id int, image tb_images) bool {
	images := new(tb_images)
	fmt.Println(image)
	t.Db.Model(&images).Where("id = ?", id).UpdateColumns(image)

	return true
}

func (t *Comic) DeleteChapter(id int) {
	var chapter tb_chapter
	t.Db.Where("id = ?", id).Delete(chapter)
}
