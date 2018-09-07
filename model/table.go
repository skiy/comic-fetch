package model

type Table struct {
	Books   tb_books
	Chapter tb_chapter
	Images  tb_images
}

type tb_books struct {
	Id     int
	Name   string
	Status int
	OriginUrl,
	OriginWeb string
	OriginBookId int
	UpdatedAt,
	CreatedAt int64
}

type tb_chapter struct {
	Id,
	Bid,
	ChapterId int
	Title string
	OrderId int
	OriginUrl string
	CreatedAt int64
}

type tb_images struct {
	Id,
	Bid,
	Cid,
	ChapterId int
	ImageUrl,
	OriginUrl string
	OrderId,
	IsRemote int
	CreatedAt int64
}
