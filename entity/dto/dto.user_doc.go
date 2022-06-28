package dto

type UserBaseData struct {
	// 使用userId作为文档的主键
	UserId   int64     `bson:"_id"`
	Age      int       `bson:"age"`
	Country  string    `bson:"country"`
	Integral int       `bson:"integral"`
	Applied  []Applied `bson:"applied"`

	// 好友列表 Friend list
	Friends []int64 `bson:"friends"`

	NickName    string `bson:"nickname"`
	AvatarURL   string `bson:"avatar_url"`
	IsOnline    bool   `bson:"online"`
	OfflineTime int64  `bson:"offline_time"`
}

type Applied struct {
	UserId int64 `bson:"user_id"`
	/* 申请状态   0 我申请别人   1  别人申请我  2  我拒绝别人  3 别人拒绝我  4 我通过别人  5 别人通过我*/
	Status int `bson:"status"`
}

const (
	Apply = iota
	OtherApply
	NoPass
	OtherNoPass
	Agree
	OtherAgree
)
