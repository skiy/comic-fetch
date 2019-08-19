package config

var (
	// WebURL 支持的网站列表
	// flag 默认使用哪种方式采集
	WebURL = map[string]map[string]string{
		"manhuaniu": {
			"pc":     "https://www.manhuaniu.com",
			"mobile": "https://m.manhuaniu.com",
			"flag":   "pc",
			"name":   "漫画牛",
		},
		"mh1234": {
			"pc":     "https://www.mh1234.com",
			"mobile": "https://www.mh1234.com",
			"flag":   "mobile",
			"name":   "漫画1234",
		},
	}
)
