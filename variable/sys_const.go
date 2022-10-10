package variable

const (
	TRUE  = "true"
	FALSE = "false"
)

const (
	// RunMode 启动模式
	RunMode = "mode"

	// RunPort 启动端口
	RunPort = "port"

	// ConfPath 配置文件
	ConfPath = "conf"

	// ReleaseMode 启动模式
	ReleaseMode = "release"

	// TimeFormat 时间格式化
	TimeFormat = "2006-01-02 15:04:05"
)

const (
	ResponseOkCode          = 400000 // 常规响应Code，所有动作，成功之后的返回码
	RespinseLimitErr        = 200503 // 超载 服务器暂时无法处理客户端的请求
	ResponseInternalErrCode = 500000 // 系统内部错误响应Code
)
