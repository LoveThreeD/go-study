package dto

type UserBaseData struct {
	UserId   int64
	Age      int
	Country  string
	Integral int
	Applied  []Applied

	// 好友列表 Friend list
	Friends []int64

	NickName    string
	AvatarURL   string
	IsOnline    bool
	OfflineTime int64
}

type Applied struct {
	UserId int64
	/* 申请状态   0 我申请别人   1  别人申请我  2  我拒绝别人  3 别人拒绝我  4 我通过别人  5 别人通过我*/
	Status int
}

const (
	Apply = iota
	OtherApply
	NoPass
	OtherNoPass
	Agree
	OtherAgree
)
