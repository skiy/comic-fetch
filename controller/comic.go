package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/skiy/comicFetch/library"
	"github.com/skiy/comicFetch/model"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Init struct {
	Model   model.Comic
	Cache   *redis.Client
	Ftimage model.Table
	Conf    library.Config
}

/**
抓取漫画初始化
*/
func (t *Init) Construct() {
	t.newBooks()
	t.getComicList()

	//采集远程图片到本地
	if t.Conf.Setting.ImageFetch {
		t.fetchImage()
	}
}

/**
检测新书
*/
func (t *Init) newBooks() {
	comicList := t.Model.GetBookList(0)

	cacheKey := "newbooks"

	var v, cacheType string
	var err error
	var str []byte

	filepath := cacheKey + ".json"

	if t.Cache != nil {
		v, err = t.Cache.Do("get", cacheKey).String()
		//空数据
		if v == "" {
			return
		}

		str = []byte(v)
		cacheType = "cache"
	} else {
		//file 方式
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			fmt.Println("Load cache filepath fail", err)
			return
		}

		str, err = ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println("Read cache filepath fail", err)
			return
		}

		if str == nil {
			return
		}

		cacheType = "file"
	}

	if err == nil {
		//fmt.Println(v)
		type m struct {
			Id   int    `json:"id"`
			Flag string `json:"flag"`
		}
		var m1 []m

		err = json.Unmarshal(str, &m1)

		if err == nil && m1 != nil {
			for _, v2 := range m1 {
				//fmt.Println(v2)
				for _, value := range comicList {
					//fmt.Println(value.OriginBookId, v2.Id, value.OriginBookId == v2.Id)
					if value.OriginBookId == v2.Id && value.OriginFlag == v2.Flag {
						goto Next
					}
				}

				//漫画160
				if v2.Flag == "mh160" {
					t.addMh160Book(v2.Id)
				}
				fmt.Printf("\n新增漫画ID(%d), 来源：<<%s>>\n", v2.Id, v2.Flag)
			Next:
			}

			//清空缓存
			if cacheType == "cache" {
				t.Cache.Do("del", cacheKey)
			} else if cacheType == "file" {
				var emptyData []byte
				err = ioutil.WriteFile(filepath, emptyData, 0755)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

/**
获取漫画列表
*/
func (t *Init) getComicList() {
	comicList := t.Model.GetBookList(0)
	//fmt.Println(comicList)

	comicTip := "\n漫画：%s (%d), %s\n"
	for _, value := range comicList {
		//fmt.Println(index, value)

		if value.Status != 0 {
			s := "暂停更新"
			if value.Status == 2 {
				s = "完结"
			}

			fmt.Printf(comicTip, value.Name, value.Id, s)
			continue
		}

		log.Printf(comicTip, value.Name, value.Id, "正在更新……")

		if strings.EqualFold(value.OriginFlag, "mh160") {
			var mh mh160
			mh.db = t.Model.Db
			mh.id = value.OriginBookId
			mh.imageUrl = value.ImageUrl
			mh.originImageUrl = value.OriginImageUrl
			mh.originPathUrl = value.OriginPathUrl
			mh.Conf = t.Conf
			mh.Init()
		}
	}
}

/**
添加漫画160的漫画
*/
func (t *Init) addMh160Book(id int) {
	var mh mh160
	mh.db = t.Model.Db
	mh.id = id
	mh.Init()
}

var (
	wg       sync.WaitGroup
	taskLoad int
	cpuNum   int
)

func (t *Init) fetchImage() {
	list := t.Model.FetchImageList()
	if len(list) == 0 {
		return
	}

	//for _, value := range list {
	//	fmt.Println(value)
	//	//break
	//}

	taskLoad = len(list)
	cpuNum = runtime.NumCPU()

	tasks := make(chan model.FtImages, taskLoad)

	fmt.Println(t.Conf.Image)

	exist, err := library.PathExists(t.Conf.Image.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !exist {
		err = os.MkdirAll(t.Conf.Image.Path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	wg.Add(cpuNum)
	for qr := 1; qr <= cpuNum; qr++ {
		go t.worker(tasks, qr)
	}

	for post := 0; post < taskLoad; post++ {
		tasks <- list[post]
	}

	close(tasks)

	wg.Wait()
}

func (t *Init) worker(tasks chan model.FtImages, worker int) {
	defer wg.Done()

	for {
		task, ok := <-tasks
		if !ok {
			fmt.Printf("Worker: %d : Shutting Down\n", worker)
			return
		}

		filename := fmt.Sprintf("%s-%s-%s", task.Bid, task.Cid, task.OrderId)

		var imagePath string
		if strings.HasPrefix(t.Conf.Image.Path, "/") {
			imagePath = t.Conf.Image.Path
		} else {
			var err error
			imagePath, err = library.GetCurrentDirectory()
			if err != nil {
				fmt.Println(err, "GetCurrentDirectory error")
				imagePath = t.Conf.Image.Path
			}
		}
		err, filename, size := library.FetchFile(task.ImageUrl, filename, imagePath, task.OriginUrl)
		if err != nil {
			fmt.Println("fetch file", err)
		} else {
			image := map[string]interface{}{"image_url": filename, "is_remote": 0, "size": size}
			tid, _ := strconv.Atoi(task.Id)
			t.Model.UpdateImageField(tid, image)
		}
	}
}
