package api

import (
	"strings"
)

// Base Base
type Base struct {
}

// sort 构造排序
func (t *Base) sort(str string, filter map[string]bool) string {
	arr := strings.Split(str, ",")
	var sortArr []string
	var order string
	for _, v := range arr {
		if s1 := strings.Trim(v, " "); s1 != "" {
			field := s1[1:]
			if s1[0] == '-' {
				order = "desc"
			} else if s1[0] == '+' {
				order = "asc"
			} else {
				order = "asc"
				field = s1
			}

			// 过滤字段
			if f, ok := filter[field]; ok && f {
				continue
			}

			if field != "" {
				sortArr = append(sortArr, field+" "+order)
			}
		}
	}

	return strings.Join(sortArr, ",")
}
