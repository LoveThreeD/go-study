package entity

type AccountData struct {
	UserID      int    `db:"userId" form:"userId"`
	Account     string `db:"account" form:"account"`
	Passwd      string `db:"passwd"  form:"passwd"`
	EquipmentID string `db:"equipmentId"  form:"equipmentId"`
}

type BaseData struct {
	UserID      int    `db:"user_id" form:"userId"`
	NickName    string `db:"nickname" form:"nickName"`
	AvatarURL   string `db:"avatar_url" form:"avatarURL"`
	Score       uint64 `db:"score" form:"score"`
	IsOnline    bool   `db:"is_online" form:"isOnline"`
	OfflineTime int64  `db:"offline_time" form:"offlineTime"`
}

type GameData struct {
	UserID   int    `db:"user_id" form:"userId"`
	GameData []byte `db:"game_data" form:"gameData"`
}

type UserBaseData struct {
	BaseData BaseData

	UserId int64 `db:"user_id"`
	Age    int   `db:"age"`
}
