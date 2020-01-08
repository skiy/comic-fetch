package command

import (
	"github.com/skiy/comic-fetch/app/model"
)

// Mh1234 漫画1234
type Mh1234 struct {
	base
}

// NewMh1234 Mh1234 init
func NewMh1234(books *model.TbBooks) *Mh1234 {
	t := &Mh1234{}
	t.Books = books
	t.ResURL = "https://mhpic.dongzaojiage.com"

	t.Prep.Book = `#Cover>img`
	t.Prep.SiteURL = "%s/wap/comic/%d.html"
	t.Prep.ChapterList = `#chapter-list-1 li a`
	t.Prep.ChapterURL = `\/comic\/[0-9]*/([0-9]*).html`
	t.Prep.ChapterPath = `chapterPath = "([^"]*)"`
	t.Prep.ImageStr = `chapterImages = \[([^\]]*)\]`
	t.Prep.Chapter = `SinMH\.initChapter\(([^;]*)\)`
	t.Prep.ImagesURL = `"([^"]*)"`
	return t
}
