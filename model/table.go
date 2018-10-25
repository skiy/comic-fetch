package model

type Table struct {
	Books   TbBooks
	Chapter TbChapter
	Images  TbImages
}

type TbBooks struct {
	Id int
	Name,
	ImageUrl string //漫画图标地址
	Status     int //状态
	OriginUrl, //漫画地址
	OriginWeb, //源站名
	OriginFlag, //源站标识
	OriginPathUrl, //上次采集图片实际路径
	OriginImageUrl string //源站漫画图标地址
	OriginBookId int //本书ID
	UpdatedAt,   //更新时间
	CreatedAt int64 //创建时间
}

type TbChapter struct {
	Id,
	Bid,
	ChapterId int
	Title string
	OrderId,
	OriginId int
	OriginUrl string
	CreatedAt int64
}

type TbImages struct {
	Id,
	Bid,
	Cid,
	ChapterId int
	ImageUrl,
	OriginUrl string
	Size int64
	OrderId,
	IsRemote int
	CreatedAt int64
}

//采集图片
type FtImages struct {
	Id,
	Bid,
	Cid,
	OrderId,
	ImageUrl,
	Name,
	OriginUrl,
	OriginFlag string
}
