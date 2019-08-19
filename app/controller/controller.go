package controller

// Controller app
type Controller interface {
	// 设置端口 (Web)
	SetPort(port int)

	// 启动服务
	Start() (err error)
}
