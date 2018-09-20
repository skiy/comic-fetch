package controller

import (
	"code.aliyun.com/skiystudy/comicFetch/library"
	"code.aliyun.com/skiystudy/comicFetch/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type mh160 struct {
	id    int
	url   string
	model model.Comic
	db    *gorm.DB
	originImageUrl,
	originWeb,
	originFlag string
}

func (t *mh160) Init() {
	t.model.Db = t.db

	t.originFlag = "mh160"
	t.originWeb = "漫画160"
	t.url = "https://m.mh160.com" //手机版

	t.mobileChapter()
}

/**
移动端获取章节
*/
func (t *mh160) mobileChapter() {
	bookUrl := t.url + fmt.Sprintf("/kanmanhua/%d/", t.id)
	fmt.Printf("正在采集漫画, URL: %s\n", bookUrl)

	doc := library.FetchSource(bookUrl)

	nowTime := time.Now().Unix()

	var bookName string
	doc.Find(".main-bar").Each(func(i int, s *goquery.Selection) {
		bookName = s.Find("h1").Text()
	})

	fmt.Printf("漫画名:《%s》\n", bookName)

	book := t.model.Table.Books
	t.model.Db.Where("name = ?", bookName).First(&book)

	if book.Id == 0 {
		books := t.model.Table.Books
		//books.Id = t.id
		books.Name = bookName
		books.Status = 0
		books.OriginUrl = bookUrl
		books.OriginWeb = t.originWeb
		books.OriginFlag = t.originFlag
		books.OriginBookId = t.id
		books.UpdatedAt = nowTime
		books.CreatedAt = nowTime

		book = t.model.CreateBook(books)
	}

	chapterList := t.model.GetChapterList(book.Id)
	var chapterIds []string
	for _, value := range chapterList {
		chapterIds = append(chapterIds, strconv.Itoa(value.OrderId))
	}
	//fmt.Println(chapterIds)

	doc.Find(".chapter-list ul li").Each(func(i int, s *goquery.Selection) { //手机版
		chapterName := s.Find("a").Text()
		url, _ := s.Find("a").Attr("href")

		//fmt.Printf("正在采集章节: %s, URL: %s \n", chapterName, t.url+url)

		var err error
		var chapterNum int

		preg := `第([0-9]*)话`
		re := regexp.MustCompile(preg)
		test := re.FindStringSubmatch(chapterName)

		if len(test) >= 2 {
			//log.Fatalf("获取章节ID失败: %s %s", url, chapterName)
			chapterNum, err = strconv.Atoi(test[1])
			if err != nil {
				log.Fatalf("章节转Int型失败: %s %s", test[1], chapterName)
			}
		}

		//preg := `2[0-9-\s:]*`
		preg = `/([0-9\/]*)/([0-9\.]*).html`
		re = regexp.MustCompile(preg)
		test = re.FindStringSubmatch(url)

		if len(test) < 3 {
			log.Fatalf("获取章节失败: %s", url)
		}

		var originChapterId int
		originChapterId, err = strconv.Atoi(test[2])
		if err != nil {
			log.Fatalf("章节ID转Int型失败: %s %s", test[1], chapterName)
		}

		has := t.InArray(test[2], chapterIds)
		if !has {
			fmt.Printf("正在采集章节: %s, URL: %s \n", chapterName, t.url+url)

			chapter := t.model.Table.Chapter
			chapter.Bid = book.Id
			chapter.ChapterId = chapterNum
			chapter.Title = chapterName
			chapter.OrderId = originChapterId
			chapter.OriginUrl = t.url + url
			chapter.CreatedAt = nowTime

			chapterInfo := t.model.CreateChapter(chapter)

			//获取共几话
			counts := t.countImage(url)

			chapterName = strings.Replace(chapterName, "-", "", -1)
			chapterName = strings.Replace(chapterName, "，", "-", -1)
			chapterName = strings.Replace(chapterName, "！", "-", -1)
			chapterName = strings.Replace(chapterName, "/", "_", -1)

			//图片
			isAdd := t.detail(test[2], book.Id, chapterInfo.Id, chapterNum, bookName, chapterName, counts)

			//isAdd = true
			if isAdd {
				tk := "https://oapi.dingtalk.com/robot/send?access_token=8eaeec8ea1c97b646e85c89e884ff1cae5e5302991088f4a8d876268ce1bd59d"
				post := `
{
     "msgtype": "text",
     "text": {
         "content": "漫画《%s》\n更新章节《%s》"
     },
     "at": {
         "atMobiles": [
             "%s"
         ], 
         "isAtAll": %s
     }
 }`
				post = fmt.Sprintf(post, bookName, chapterName, "18565756628", "false")
				//fmt.Println(post)
				req, err := http.NewRequest("POST", tk, strings.NewReader(post))
				if err != nil {
					fmt.Println(err)
				}

				req.Header.Add("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
				}

				defer resp.Body.Close()
				//
				//body, err := ioutil.ReadAll(resp.Body)
				//if err != nil {
				//	fmt.Println(err)
				//}
				//
				//fmt.Println(string(body))
			}
		}
	})
}

/**
获取共几话
*/
func (t *mh160) countImage(url string) (counts int) {
	fetchUrl := t.url + url
	//fmt.Println(fetchUrl)

	var err error
	doc := library.FetchSource(fetchUrl)

	doc.Find(".main-bar").Each(func(i int, s *goquery.Selection) {
		imagePage := s.Find(".manga-page").Text()

		preg := `([0-9\/]*)/([0-9\/]*)P`
		re := regexp.MustCompile(preg)
		test := re.FindStringSubmatch(imagePage)

		if len(test) < 3 {
			log.Fatalf("获取章节话数失败: %s %s", url, imagePage)
		}

		counts, err = strconv.Atoi(test[2])
		if err != nil {
			fmt.Println(err)
		}
	})

	return counts
}

/**
获取漫画图片
*/
func (t *mh160) detail(originChapterId string, bookId, chapterId, chapterNum int, bookName, chapterName string, counts int) (isAdd bool) {
	var realUrl string
	var has bool

	baseUrl := "https://mhpic%s.lineinfo.cn/mh160tuku/%s/%s_%d/%s_%s/"
	//fmt.Println("originImageUrl", t.originImageUrl)

	//有源
	if t.originImageUrl != "" {
		preg := `https:\/\/mhpic([5-7])\.lineinfo\.cn\/mh160tuku\/([a-z]*)\/([^_]*)_([0-9]*)\/([^_]*)_([0-9]*)\/00([0-9]*)\.jpg`
		reg := regexp.MustCompile(preg)
		test := reg.FindStringSubmatch(t.originImageUrl)
		//fmt.Println(test, len(test))

		if len(test) == 8 {
			realUrl = fmt.Sprintf(baseUrl, test[1], test[2], test[3], t.id, chapterName, originChapterId) + "00%s.jpg"
		}
	} else {
		realUrl = t.getImageUrl(baseUrl, bookName, chapterName, originChapterId, bookId)

		//fmt.Println(realUrl)
		if realUrl != "" {
			has = true
		}
		//fmt.Println("originImageUrl", t.originImageUrl)
	}

	if t.originImageUrl == "" {
		t.model.DeleteChapter(chapterId)
		fmt.Println("该话漫画暂时获取不到")
		isAdd = false
		return
	}

	//fmt.Println(realUrl)

	images := t.model.Table.Images
	images.Bid = bookId
	images.Cid = chapterId
	images.ChapterId = chapterNum
	images.ImageUrl = ""
	images.IsRemote = 1
	images.CreatedAt = time.Now().Unix()

	var fix string
	for i := 1; i < counts; i++ {
		//fmt.Println(i)
		if i < 10 {
			fix = "0" + strconv.Itoa(i)
		} else {
			fix = strconv.Itoa(i)
		}

		images.OriginUrl = strings.Replace(fmt.Sprintf(realUrl, fix), " ", "", -1)

		images.OrderId = i

		if !has && i == 1 {
			refererUrl := fmt.Sprintf("/kanmanhua/%d/%s.html", t.id, originChapterId)
			isRight := t.checkUrl(images.OriginUrl, refererUrl)
			if !isRight {
				t.originImageUrl = ""

				realUrl = t.getImageUrl(baseUrl, bookName, chapterName, originChapterId, bookId)

				if realUrl == "" {
					t.model.DeleteChapter(chapterId)
					fmt.Println("该话漫画暂时获取不到")
					isAdd = false
					return
				}

				images.OriginUrl = strings.Replace(fmt.Sprintf(realUrl, fix), " ", "", -1)
			}
		}

		t.model.CreateImages(images)
	}

	isAdd = true
	return

}

func (t *mh160) getImageUrl(baseUrl, bookName, chapterName, originChapterId string, bookId int) (realUrl string) {
	imageUrl := fmt.Sprintf(baseUrl, "%d", "%s", bookName, t.id, chapterName, originChapterId)
	pathUrl := imageUrl + "0001.jpg"
	chapterUrl := fmt.Sprintf("/kanmanhua/%d/%s.html", t.id, originChapterId)
	//fmt.Println(chapterUrl)

	var mhpic = [...]int{5, 6, 7}
	var pathUrl2 string

	for _, picNum := range mhpic {
		for i := 122; i >= 97; i-- {
			c := string(i)
			//fmt.Println(picNum, c)
			pathUrl2 = strings.Replace(fmt.Sprintf(pathUrl, picNum, c), " ", "", -1)

			//fmt.Println(pathUrl2)

			isRight := t.checkUrl(pathUrl2, chapterUrl)
			if isRight {
				//fmt.Println(pathUrl2)
				realUrl = strings.Replace(pathUrl2, "01.jpg", "%s.jpg", -1)
				fmt.Printf("当前漫画的 PATH 是: %s\n", pathUrl2)
				t.originImageUrl = pathUrl2

				t.model.UpdateBookImageUrl(bookId, t.originImageUrl)
				break
			}
		}

		if t.originImageUrl != "" {
			break
		}
	}

	if realUrl == "" {
		fmt.Printf("获取漫画图片失败,此URL: %s\n", pathUrl2)
	}

	return
}

func (t *mh160) checkUrl(url, chapterUrl string) bool {
	//str := "https://mhpic6.lineinfo.cn/mh160tuku/d/斗罗大陆2绝世唐门_11140/第82话极动中的炽烈—天帝之锤_488477/0001.jpg"
	//fmt.Println(url, "\n", str)
	client := &http.Client{}

	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//增加header选项
	reqest.Header.Add("NT", "1")
	reqest.Header.Add("If-Modified-Since", "Thu, 06 Sep 2018 03:54:19 GMT")
	reqest.Header.Add("If-None-Match", "BDE9E8B0317BF99A37BE8FE52763AF1E")
	reqest.Header.Add("Referer", t.url+chapterUrl)

	//处理返回结果
	res, _ := client.Do(reqest)
	defer res.Body.Close()

	//fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return false
	}

	return true
}

func (t *mh160) InArray(str string, arr []string) bool {
	for _, value := range arr {
		if value == str {
			return true
		}
	}
	return false
}
