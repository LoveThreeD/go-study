package dto

type UserBaseData struct {
	//BaseData entity.BaseData
	UserId   int64
	Age      int
	Country  string
	Integral int
	// XXX 以下在mongo中有无更合适的存储方式 性能呢: (改变： 存储为一个结构体, 从搜索速度上选择4个数组，但是好友业务中一个功能发起多条mongo数组处理,为结构体时只需一条，性能可能更优)
	/*// 申请列表(别人给我的申请) Application list
	ApplicationList []int64
	// 已申请列表(我给别人的申请) AlreadyApplied List
	AlreadyAppliedList []int64
	// 已拒绝列表(我拒绝别人) NoPass List
	NoPassList []int64
	// 被拒绝列表(别人拒绝我) NoPassRecord list
	NoPassRecord []int64*/

	Applied []Applied

	// 好友列表 Friend list
	Friends []int64

	NickName    string
	AvatarURL   string
	IsOnline    bool
	OfflineTime int64
}

type Applied struct {
	UserId int64
	/*申请状态   0 我申请别人   1  别人申请我  2  我拒绝别人  3 别人拒绝我  4 我通过别人  5 别人通过我*/
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
