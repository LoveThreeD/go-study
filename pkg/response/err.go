package response

// 错误码规则:
// (1) 错误码需为 > 0 的数;
//
// (2) 错误码为 5 位数:
//              ----------------------------------------------------------
//                  第1位               2、3位                  4、5位
//              ----------------------------------------------------------
//                服务级错误码          模块级错误码	         具体错误码
//              ----------------------------------------------------------

var (
	MsgConventErr = "类型转换异常"

	MsgMongoDbConnectionError   = "Mongo数据库连接失败"
	MsgMongoCollConnectionError = "Mongo文档连接失败"
	MsgMongoCreateUserError     = "Mongo创建用户失败"
	MsgMongoSelectUserError     = "Mongo查询用户失败"

	MsgFailed        = "处理失败"
	MsgInitDataError = "初始化数据失败"
)

var (
	OK  = response(200, "ok")   // 通用成功
	Err = response(500, "调用失败") // 通用错误

	// ErrParam 服务级错误码
	ErrParam           = response(10001, "参数有误")
	ErrSignParam       = response(10002, "签名参数有误")
	ErrNoMethodHandler = response(10003, "函数未发现")
	ErrNoRouteHandler  = response(10004, "路由未发现")
	ErrNoPermission    = response(10005, "没有权限")

	// ErrUserService 模块级错误码 - 用户模块
	ErrUserService        = response(20100, "用户服务异常")
	ErrAccountAndPassword = response(20101, "账户或密码不正确")
)
