package model

import (
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
	t.Db.Table("tb_images AS i").Select("i.bid,i.cid,i.order_id,i.origin_url AS image_url,b.name,b.origin_url,b.origin_flag").Joins("LEFT JOIN tb_books b ON b.id = i.bid").Where("I.image_url = ''").Order("i.id ASC").Scan(&images)
	return
}

func (t *Comic) GetChapterList(bid int) (chapters []tb_chapter) {
	if bid != 0 {
		t.Db.Where("bid = ?", bid)
	}
	t.Db.Find(&chapters)
	return
}

func (t *Comic) UpdateBookImageUrl(id int, url string) bool {
	books := new(tb_books)
	//fmt.Println(id, url)
	t.Db.Model(&books).Where("id = ?", id).UpdateColumn("origin_image_url", url)

	return true
}

func (t *Comic) DeleteChapter(id int) {
	var chapter tb_chapter
	t.Db.Where("id = ?", id).Delete(chapter)
}
