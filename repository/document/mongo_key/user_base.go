package mongo_key

const (
	BaseUserId            = "_id"
	BaseAge               = "age"
	BaseCountry           = "country"
	BasePoints            = "points"
	BaseApplied           = "applied"
	BaseUserIdInApplied   = "applied.user_id"
	BaseUpUserIdInApplied = "applied.$.user_id"
	BaseStatusInApplied   = "applied.status"
	BaseUpStatusInApplied = "applied.$.status"
	BaseFriends           = "friends"
	BaseNickName          = "nickname"
	BaseAvatarUrl         = "avatar_url"
	BaseIsOnline          = "online"
	BaseOfflineTime       = "offline_time"
)
