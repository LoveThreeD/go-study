package mongo_key

const (
	BaseUserId            = "_id"
	BaseAge               = "age"
	BaseCountry           = "country"
	BasePoints            = "points"
	BaseApplied           = "request_status"
	BaseUserIdInApplied   = "request_status.user_id"
	BaseUpUserIdInApplied = "request_status.$.user_id"
	BaseStatusInApplied   = "request_status.status"
	BaseUpStatusInApplied = "request_status.$.status"
	BaseFriends           = "friends"
	BaseNickName          = "nickname"
	BaseAvatarUrl         = "avatar_url"
	BaseIsOnline          = "online"
	BaseOfflineTime       = "offline_time"
)
