package command

import (
	"github.com/skiy/comic-fetch/app/model"
)

// Nfydd 漫画牛
type Nfydd struct {
	base
}

// NewNfydd Nfydd init
func NewNfydd(books *model.TbBooks) *Nfydd {
	t := &Nfydd{}
	t.Books = books
	t.ResURL = "https://res.nbhbzl.com"

	t.Prep.Book = `img.pic`
	t.Prep.SiteURL = "%s/manhua/%d/"
	t.Prep.ChapterList = `#chapter-list-1 li a`
	t.Prep.ChapterURL = `\/([0-9]*).html`
	t.Prep.ChapterPath = `chapterPath = "([^"]*)"`
	t.Prep.ImageStr = `chapterImages = \[([^\]]*)\]`
	t.Prep.Chapter = `SinMH\.initChapter\(([^;]*)\)`
	t.Prep.ImagesURL = `"([^"]*)"`
	return t
}
