package cerror

import "net/http"

const (
	// ErrSuccess 成功 200
	ErrSuccess = iota
	// ErrFailure 失败
	ErrFailure
	// ErrAddSuccess 创建成功 201
	ErrAddSuccess
	// ErrInvalidParameter 参数错误
	ErrInvalidParameter
	// ErrUnauthorized 未授权 401
	ErrUnauthorized
	// ErrForbidden 权限不足 403
	ErrForbidden

	// ErrEmptyUsernameParam 用户名不存在
	ErrEmptyUsernameParam
	// ErrEmptyPasswordParam 密码不存在
	ErrEmptyPasswordParam
	// ErrUserDisabled 用户已禁用
	ErrUserDisabled
	// ErrPasswordIncorrect 密码错误
	ErrPasswordIncorrect
	// ErrPasswordRepeat 密码不一敌
	ErrPasswordRepeat

	// ErrNoUser 用户不存在
	ErrNoUser
	// ErrUserExist 用户已存在
	ErrUserExist
	// ErrLoginSucceeded 登录成功
	ErrLoginSucceeded
	// ErrLoginFailed 登录失败
	ErrLoginFailed

	// ErrNoData 无数据
	ErrNoData

	// ErrDeviceOffline 设备不在线
	ErrDeviceOffline
	// ErrAllDeviceOffline 所有设备不在线
	ErrAllDeviceOffline
	// ErrSendToDeviceSuccess 下发消息给设备成功
	ErrSendToDeviceSuccess
	// ErrSendToDeviceFailure 下发消息给设备失败
	ErrSendToDeviceFailure
	// ErrGetData 数据获取失败 (select)
	ErrGetData
	// ErrAddData 数据创建失败 (create)
	ErrAddData
	// ErrUpdateData 数据更新失败 (update)
	ErrUpdateData
	// ErrDeleteData 数据删除失败 (delete)
	ErrDeleteData
	// ErrBookNameNotExist 漫画名不存在
	ErrBookNameNotExist
	// ErrBookIDNotExist 漫画 ID 不存在
	ErrBookIDNotExist
	// ErrChapterIDNotExist 漫画章节 ID 不存在
	ErrChapterIDNotExist
	// ErrImageIDNotExist 漫画图库不存在
	ErrImageIDNotExist
)

// ErrMessage 错误信息
var ErrMessage = map[int]string{
	ErrSuccess:          "操作成功",
	ErrFailure:          "操作失败",
	ErrAddSuccess:       "创建成功",
	ErrInvalidParameter: "参数错误",
	ErrUnauthorized:     "未授权",
	ErrForbidden:        "权限不足",

	ErrEmptyUsernameParam: "用户名参数为空",
	ErrEmptyPasswordParam: "密码参数为空",
	ErrUserDisabled:       "账号已被禁用",
	ErrPasswordIncorrect:  "密码不正确",
	ErrPasswordRepeat:     "两次输入密码不一致",
	ErrGetData:            "数据获取失败",
	ErrAddData:            "数据创建失败",
	ErrUpdateData:         "数据更新失败",
	ErrDeleteData:         "数据删除失败",

	ErrUserExist: "用户名已存在",
	ErrNoUser:    "用户不存在",

	ErrLoginSucceeded: "登录成功",
	ErrLoginFailed:    "登录失败",

	ErrNoData: "数据不存在",

	ErrDeviceOffline:       "设备离线",
	ErrAllDeviceOffline:    "所有设备离线",
	ErrSendToDeviceSuccess: "操作下发成功",
	ErrSendToDeviceFailure: "操作下发失败",

	ErrBookNameNotExist:  "漫画名不存在",
	ErrBookIDNotExist:    "漫画 ID 不存在",
	ErrChapterIDNotExist: "漫画章节 ID 不存在",
	ErrImageIDNotExist:   "漫画图库不存在",
}

// GetMessage get Err message
func GetMessage(flag int) string {
	if ErrMessage[flag] != "" {
		return ErrMessage[flag]
	}

	return ""
}

// GetHTTPStatus http status code
func GetHTTPStatus(flag int) int {
	var status int

	switch flag {
	case ErrSuccess:
		status = http.StatusOK
		break

	case ErrUnauthorized:
		status = http.StatusUnauthorized
		break

	case ErrForbidden:
		status = http.StatusForbidden
		break

	case ErrAddSuccess:
		status = http.StatusCreated
		break

	default:
		status = http.StatusBadRequest

	}

	return status
}

// GetErrMsg 获取错误信息及 http status
func GetErrMsg(flag int) (string, int) {
	return GetMessage(flag), GetHTTPStatus(flag)
}
