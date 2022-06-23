package friend_dto

type ReqFriendSearch struct {
	NickName string `json:"nickName" form:"nickName" bson:"nickname"`
}
