package lnotify

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/util/gconv"
	"github.com/skiy/comic-fetch/app/library/lcfg"
)

// NotifyMessage 通知消息
type NotifyMessage struct {
	// 通知类型: 1. 新增, 2. 更新
	flag int
}

// NewNotifyMessage NotifyMessage init
func NewNotifyMessage(flag int) *NotifyMessage {
	t := &NotifyMessage{}
	t.flag = flag
	return t
}

// Dingtalk 钉钉通知
func (t *NotifyMessage) Dingtalk(params ...interface{}) (err error) {
	cfg := lcfg.GetCfg()

	notifyURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", cfg.GetString("notify.dingtalk.robot_access_token"))

	var title, text, picURL, messageURL string
	picURL = gconv.String(params[2])
	messageURL = gconv.String(params[3])

	// 新增
	if t.flag == 1 {
		text = fmt.Sprintf("来源: %s", params[1])
		title = fmt.Sprintf("新增漫画:《%s》", params[0])
	} else if t.flag == 2 { // 更新
		text = fmt.Sprintf("章节:《%s》", params[1])
		title = fmt.Sprintf("更新漫画:《%s》", params[0])
	}

	postParams := g.Map{
		"msgtype": "link",
		"link": g.Map{
			"text":       text,
			"title":      title,
			"picUrl":     picURL,
			"messageUrl": messageURL,
		},
	}

	c := ghttp.NewClient()
	c.SetHeader("Content-Type", "application/json")

	paramsJSON := gconv.String(postParams)

	r, err := c.Post(notifyURL, paramsJSON)
	if err != nil {
		return fmt.Errorf("notify paramsJsonStr: %v, err: %s", params, err.Error())
	}
	defer r.Close()

	// 提交成功
	if r.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("notify http response status code: %d", r.StatusCode)
}
