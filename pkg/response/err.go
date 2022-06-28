package response

import "net/http"

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
	/* mongo学习 */
	MsgMongoDbConnectionError   = "Mongo数据库连接失败"
	MsgMongoCollConnectionError = "Mongo文档连接失败"
	MsgMongoCreateUserError     = "Mongo创建用户失败"
	MsgMongoSelectUserError     = "Mongo查询用户失败"
	MsgMongoUpdateUserError     = "Mongo更新用户失败"

	MsgFailed        = "处理失败"
	MsgInitDataError = "初始化数据失败"
	/* 好友 */
	MsgFriendNumberError = "好友数量超出"
)

var (
	OK = response(http.StatusOK, "ok") // 通用成功
)
