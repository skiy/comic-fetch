package library

import (
	"fmt"
	"net/http"
	"strings"
)

type Message struct {
	Conf Config
}

/**
钉钉通知
flag 1新增漫画, 2更新漫画
*/
func (t *Message) Dingtalk(flag int, params ...string) (notice bool) {
	tk := "https://oapi.dingtalk.com/robot/send?access_token=" + t.Conf.Notice.DindTalkToken

	var content, mobile, all string

	mobile = t.Conf.Notice.Receiver
	all = "false"

	if t.Conf.Notice.Receiver == "" {
		all = "true"
	}

	post := `
{
     "msgtype": "text",
     "text": {
         "content": "%s"
     },
     "at": {
         "atMobiles": [
             "%s"
         ], 
         "isAtAll": %s
     }
 }`
	//新增
	if flag == 1 {
		content = fmt.Sprintf("新增\n漫画:《%s》\n来源: %s\n", params[0], params[1])
		//更新
	} else if flag == 2 {
		content = fmt.Sprintf("更新\n漫画:《%s》\n章节:《%s》\n", params[0], params[1])
	}

	post = fmt.Sprintf(post, content, mobile, all)
	//fmt.Println(post)
	req, err := http.NewRequest("POST", tk, strings.NewReader(post))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	notice = true
	return
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(string(body))
}
