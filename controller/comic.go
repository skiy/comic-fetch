package controller

import (
	"code.aliyun.com/skiystudy/comicFetch/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"os"
	"io/ioutil"
	"sync"
	"runtime"
	"code.aliyun.com/skiystudy/comicFetch/library"
	"crypto/md5"
	"net/http"
	"io"
)

type Init struct {
	Model model.Comic
	Cache *redis.Client
	Ftimage model.Table
	Conf library.Config
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
			fmt.Println(1, err)
			return
		}

		str, err = ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println(2, err)
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
				fmt.Printf("新增漫画ID(%d), 来源：<<%s>>\n", v2.Id, v2.Flag)
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

		fmt.Printf(comicTip, value.Name, value.Id, "正在更新……")

		if strings.EqualFold(value.OriginFlag, "mh160") {
			var mh mh160
			mh.db = t.Model.Db
			mh.id = value.OriginBookId
			mh.imageUrl = value.ImageUrl
			mh.originImageUrl = value.OriginImageUrl
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
	wg sync.WaitGroup
 	taskLoad int
 	cpuNum int
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
	fmt.Println(taskLoad)
	cpuNum = runtime.NumCPU()

	fmt.Println(taskLoad, cpuNum)

	tasks := make(chan model.FtImages, taskLoad)

	fmt.Println(t.Conf.Image)

	exist, err := library.PathExists(t.Conf.Image.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ! exist {
		err = os.MkdirAll(t.Conf.Image.Path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	wg.Add(cpuNum)
	for qr := 1; qr <= cpuNum; qr ++ {
		go t.worker(tasks, qr)
	}

	for post := 0; post < taskLoad; post++ {
		tasks <- list[post]
	}

	close(tasks)

	wg.Wait()
}

func (t *Init) worker (tasks chan model.FtImages, worker int)  {
	defer wg.Done()

	for {
		task, ok := <-tasks
		if !ok {
			fmt.Printf("Worker: %d : Shutting Down\n", worker)
			return
		}

		i := strings.LastIndex(task.ImageUrl,".")
		suffix := task.ImageUrl[i:]

		//filename, _ := fmt.Printf("%s-%s-%s %s", task.Bid, task.Cid, task.OrderId, suffix)
		filename := fmt.Sprintf("%s-%s-%s", task.Bid, task.Cid, task.OrderId)
		filenameBype := []byte(filename)
		md5Filename := md5.Sum(filenameBype)
		filename = fmt.Sprintf("%x%s", md5Filename, suffix)
		fmt.Println(filename, suffix)
		filepath := t.Conf.Image.Path + "/" + filename

		exist, err := library.PathExists(filepath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if exist {
			continue
		}

		imageFile,err := os.Create(filepath)
		if err != nil {
			fmt.Printf("[writeImage create file]: fileName: %s\n href: %s\nerror: %s\n", filepath, task.ImageUrl, err.Error())
			continue
		}

		client := &http.Client{}

		//提交请求
		reqest, err := http.NewRequest("GET", task.ImageUrl, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//增加header选项
		reqest.Header.Add("NT", "1")
		reqest.Header.Add("If-Modified-Since", "Thu, 06 Sep 2018 03:54:19 GMT")
		reqest.Header.Add("If-None-Match", "BDE9E8B0317BF99A37BE8FE52763AF1E")
		reqest.Header.Add("Referer", task.OriginUrl)

		//处理返回结果
		res, _ := client.Do(reqest)

		//fmt.Println(res.StatusCode)
		if res.StatusCode != 200 {
			fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
			os.Remove(filepath)
			continue
		}

		size,err := io.Copy(imageFile, res.Body)
		if err != nil {
			fmt.Printf("io.Copy: error: %s  href: %s\n", err.Error(), task.ImageUrl)
			os.Remove(filepath)
		}
		fmt.Printf("Get From %s: %d bytes\n", task.ImageUrl, size)

		//fmt.Println(task, suffix)
	}
}
