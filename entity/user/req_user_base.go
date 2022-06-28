package user

type ReqUserBase struct {
	UserId   int64  `form:"userId"`
	NickName string `form:"nickName"`
	Country  string `form:"country"`
}
