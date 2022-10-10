package result

// 数据操作类Code
const (
	SelectErrCode          = 400400 // 数据查询失败
	SelectEmptyErrCode     = 400401 // 数据查询为空
	SelectExistErrCode     = 400402 // 数据已存在
	UpdateErrCode          = 400403 // 数据更新失败
	DeleteErrCode          = 400404 // 数据删除失败
	CreateErrCode          = 400405 // 数据存储失败
	ReqSendErrCode         = 400406 // 请求失败
	ExportModelDataErrCode = 400407 // 导出模型错误
	ImportModelDataErrCode = 400408 // 导入模型错误
)

// 参数校验类Code
const (
	ParamCheckErrCode = 400200 // 参数检验失败
	ParamParseErrCode = 400202 // 参数解析失败
)

// AuthFailCode 用户响应类Code
const (
	AuthFailCode = 400100 // 认证失败
	AuthErrCode  = 400110 // 租户用户权限验证失败
)

var (
	AuthFailedErr  = &ConstErr{ErrStr: "认证失败", ErrCode: AuthFailCode}
	AuditFailedErr = &ConstErr{ErrStr: "审计信息落库失败", ErrCode: AuthFailCode}
)
