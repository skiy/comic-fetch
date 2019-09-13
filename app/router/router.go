package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/skiy/comic-fetch/app/controller/web"
	"github.com/skiy/comic-fetch/app/controller/web/api"
)

func init() {
	s := g.Server()

	// Static setting
	//s.SetIndexFolder(true)
	//s.SetServerRoot(".")

	//s.AddStaticPath("/", "home/index.html")
	//s.AddStaticPath("/static", "home/static")

	webHome := web.NewHomeController()
	s.BindHandler("/", webHome.Index)

	apiGroup := s.Group("/api")

	apiBooks := api.NewBook()
	apiGroup.GET("/books", apiBooks.List)
	apiGroup.GET("/books/:id", apiBooks.List)
	apiGroup.POST("/books", apiBooks.Add)
	apiGroup.PUT("/books/:id", apiBooks.Update)
	apiGroup.DELETE("/books/:id", apiBooks.Delete)

	apiChapters := api.NewChapter()
	apiGroup.GET("/books/:book_id/chapters", apiChapters.List)
	apiGroup.GET("/books/:book_id/chapters/:id", apiChapters.List)
	apiGroup.GET("/books/:book_id/parts", apiChapters.List)
	apiGroup.GET("/books/:book_id/parts/:chapter_num", apiChapters.List)
	apiGroup.PUT("/books/:book_id/chapters/:id", apiChapters.Update)
	apiGroup.DELETE("/books/:book_id/chapters/:id", apiChapters.Delete)

	apiComics := api.NewComic()
	apiGroup.GET("/books/:book_id/chapters/:chapter_id/comics", apiComics.List)
	apiGroup.GET("/books/:book_id/chapters/:chapter_id/comics/:id", apiComics.List)
	apiGroup.GET("/books/:book_id/chapters/:chapter_id/parts", apiComics.List)
	apiGroup.GET("/books/:book_id/chapters/:chapter_id/parts/:comic_num", apiComics.List)
	apiGroup.DELETE("/books/:book_id/chapters/:chapter_id/comics/:id", apiComics.Delete)

	apiGroup.GET("/search/:name", apiBooks.Search)
}
