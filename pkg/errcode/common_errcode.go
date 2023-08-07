package errcode

var (
	Success                    = NewErrorCode("00000", "成功")
	ClientError                = NewErrorCode("A0001", "用户端错误")
	NotFound                   = NewErrorCode("A0100", "不存在")
	AccessPermissionError      = NewErrorCode("A0200", "访问权限异常")
	AbnormalUserOperation      = NewErrorCode("A0300", "用户操作异常")
	InvalidParams              = NewErrorCode("A0400", "无效的用户输入")
	BlankRequiredParams        = NewErrorCode("A0401", "请求必填参数为空")
	ExceedRangeParams          = NewErrorCode("A0402", "请求参数值超出允许的范围")
	MalformedParams            = NewErrorCode("A0403", "参数格式不匹配")
	ErrDeserializedParams      = NewErrorCode("A0404", "请求参数反序列化失败")
	SensitiveWordsParams       = NewErrorCode("A0405", "请求参数包含违禁敏感词")
	ServerError                = NewErrorCode("A0500", "用户请求服务异常")
	TooManyRequests            = NewErrorCode("A0501", "请求次数超出限制")
	ConcurrentExceedLimit      = NewErrorCode("A0502", "请求并发数超出限制")
	WaitUserOperation          = NewErrorCode("A0503", "用户操作请等待")
	RepeatedRequest            = NewErrorCode("A0504", "用户重复请求")
	AbnormalUserResources      = NewErrorCode("A0600", "用户资源异常")
	AbnormalUserVersion        = NewErrorCode("A0700", "用户当前版本异常")
	MismatchUserVersion        = NewErrorCode("A0701", "用户安装版本与系统不匹配")
	TooLowUserVersion          = NewErrorCode("A0702", "用户安装版本过低")
	TooHighUserVersion         = NewErrorCode("A0703", "用户安装版本过高")
	ExpiredUserVersion         = NewErrorCode("A0704", "用户安装版本已过期")
	MismatchAPIVersion         = NewErrorCode("A0705", "用户API请求版本不匹配")
	TooLowAPIVersion           = NewErrorCode("A0706", "用户API请求版本过低")
	TooHighAPIVersion          = NewErrorCode("A0707", "用户API请求版本过高")
	InternalError              = NewErrorCode("B0001", "系统执行出错")
	InvalidStartupParams       = NewErrorCode("B0002", "系统启动参数错误")
	SystemTimeout              = NewErrorCode("B0100", "系统执行超时")
	SystemResourceError        = NewErrorCode("B0200", "系统资源异常")
	ReadDiskFailed             = NewErrorCode("B0201", "系统读取磁盘文件失败")
	ThirdPartyServiceError     = NewErrorCode("C0001", "调用第三方服务出错")
	MiddlewareServiceError     = NewErrorCode("C0100", "中间件服务出错")
	MiddlewareServiceTimeout   = NewErrorCode("C0101", "中间件执行超时")
	DatabaseServiceError       = NewErrorCode("C0200", "数据库服务出错")
	DatabaseServiceTimeout     = NewErrorCode("C0201", "数据库服务超时")
	NotificationServiceError   = NewErrorCode("C0300", "通知服务出错")
	NotificationServiceTimeout = NewErrorCode("C0301", "通知服务超时")
)
