package response

import "net/http"

var (
	MsgConventErr = "type conversion fail"
	/* mongo */

	MsgMongoDbConnectionError   = "MongoDb connect fail"
	MsgMongoCollConnectionError = "Mongo document connect fail"
	MsgMongoCreateUserError     = "Mongo create user fail"
	MsgMongoSelectUserError     = "Mongo select user fail"
	MsgMongoUpdateUserError     = "Mongo update user fail"
	/* 通用*/

	MsgFailed            = "deal fail"
	MsgParamsError       = "param error"
	MsgInitDataError     = "init data fail"
	MsgFriendNumberError = "number of friends fail"

	/* 任务和关卡相关*/

	MsgPreviousError      = "pass the previous level first to continue the challenge"
	MsgNotSubsequentError = "subsequent levels are not yet open, so stay tuned"
	MsgTaskNotFoundError  = "the task is not in the task list"
	MsgTaskRepeatError    = "the task has been completed, no need to repeat it"
	MsgLevelChooseError   = "level choose error"
	MsgLevelNotFoundError = "the task is not in the task list"
	MsgLevelNotSuccess    = "level has not been completed"
)

var (
	OK = response(http.StatusOK, http.StatusText(http.StatusOK)) // 通用成功
)
