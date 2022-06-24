package friend_dto

type ReqFriendSearch struct {
	NickName string `json:"nickName" form:"nickName"`
}

type ReqFriendAdd struct {
	SelfUserId   int64 `json:"SelfUserId" form:"SelfUserId" commit:"申请人ID,也就是自身Id"`
	FriendUserId int64 `json:"friendUserId" form:"friendUserId" commit:"好友Id"`
	Agree        int   `json:"agree" form:"agree" commit:"同意还是拒绝"`
}

type ReqRecommend struct {
	UserId int64 `json:"userId" form:"userId"`
}
