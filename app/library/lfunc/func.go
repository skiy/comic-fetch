package lfunc

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/skiy/comic-fetch/app/config/cerror"
	"net/http"
)

// Response 生成 Response
// params[0] data
func Response(code int, params ...gdb.Map) (int, gdb.Map) {
	message, status := cerror.GetErrMsg(code)

	paramLen := len(params)

	if status == http.StatusCreated {
		code = 0
	}

	resp := gdb.Map{
		"code":    code,
		"message": message,
		"data":    gdb.Map{},
	}

	// 主层数据合并
	if paramLen > 0 {
		for k, v := range params[0] {
			resp[k] = v
		}
	}

	return status, resp
}
