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

/**
创建漫画
*/
func (t *Comic) CreateBook(books TbBooks) (book TbBooks) {
	t.Db.FirstOrCreate(&book, books)
	return book
}

/**
创建漫画章节
*/
func (t *Comic) CreateChapter(chapters TbChapter) (chapter TbChapter) {
	t.Db.FirstOrCreate(&chapter, chapters)
	return chapter
}

/**
创建漫画图片
*/
func (t *Comic) CreateImages(images TbImages) (image TbImages) {
	t.Db.FirstOrCreate(&image, images)
	return image
}

/**
获取漫画列表
*/
func (t *Comic) GetBookList(id int) (books []TbBooks) {
	if id != 0 {
		t.Db.Where("id = ?", id).Find(&books)
		return
	}
	t.Db.Find(&books)
	return
}

/**
获取非本地的漫画图片列表
*/
func (t *Comic) FetchImageList() (images []FtImages) {
	//sql := "SELECT i.bid,i.cid,i.order_id,i.origin_url,b.name,b.origin_url,b.origin_flag FROM tb_images AS i LEFT JOIN tb_books b ON b.id = i.bid WHERE I.image_url = '' ORDER BY i.id ASC" .Limit(10)
	t.Db.Table("tb_images AS i").Select("i.id,i.bid,i.cid,i.order_id,i.origin_url AS image_url,b.name,b.origin_url,b.origin_flag").Joins("LEFT JOIN tb_books b ON b.id = i.bid").Where("i.image_url = ''").Order("i.id ASC").Scan(&images)
	return
}

/**
漫画章节列表
*/
func (t *Comic) GetChapterList(bid int) (chapters []TbChapter) {
	if bid != 0 {
		t.Db.Where("bid = ?", bid).Order("order_id ASC").Find(&chapters)
		return
	}
	t.Db.Find(&chapters)
	return
}

/**
更新漫画
*/
func (t *Comic) UpdateBook(id int, book TbBooks) bool {
	books := new(TbBooks)
	t.Db.Model(&books).Where("id = ?", id).UpdateColumns(book)

	return true
}

/**
更新图片
*/
func (t *Comic) UpdateImage(id int, image TbImages) bool {
	images := new(TbImages)
	t.Db.Model(&images).Where("id = ?", id).UpdateColumns(image)

	return true
}

/**
更新漫画图片
*/
func (t *Comic) UpdateImageField(id int, image map[string]interface{}) bool {
	images := new(TbImages)
	t.Db.Model(&images).Where("id = ?", id).UpdateColumns(image)

	return true
}

/**
删除章节
*/
func (t *Comic) DeleteChapter(id int) {
	var chapter TbChapter
	t.Db.Where("id = ?", id).Delete(chapter)
}

/**
获取漫画图片列表
*/
func (t *Comic) GetImages(bid, cid int) (images []TbImages) {
	t.Db.Where("bid = ?", bid).Where("cid = ?", cid).Find(&images)
	return
}
