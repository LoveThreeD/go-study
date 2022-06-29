package service

import (
	"github.com/pkg/errors"
	"study/entity/dto"
	"study/entity/friend_dto"
	"study/pkg/response"
	"study/repository/document"
)

func SearchUser(search *friend_dto.ReqFriendSearch) ([]*friend_dto.RespFriendRecommend, error) {
	users, err := document.SelectUserByNickname(search)
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}
	return users, nil
}

//  RequestFriend  添加到申请列表  添加到已申请列表
func RequestFriend(friend *friend_dto.ReqFriendAdd) error {

	item := dto.Applied{
		UserId: friend.FriendUserId,
		Status: dto.Apply,
	}
	if err := document.AddApplied(friend.SelfUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	item = dto.Applied{
		UserId: friend.SelfUserId,
		Status: dto.OtherApply,
	}
	if err := document.AddApplied(friend.FriendUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	return nil
}

// AddFriend  添加到好友列表
func AddFriend(friend *friend_dto.ReqFriendAdd) error {
	// 好友限制30个
	selfData, err := document.SelectUser(int(friend.SelfUserId))
	if err != nil {
		return err
	}
	friendData, err := document.SelectUser(int(friend.FriendUserId))
	if err != nil {
		return err
	}
	if len(selfData.Friends)+1 > 30 && len(friendData.Friends)+1 > 30 {
		return errors.New(response.MsgFriendNumberError)
	}
	// XXX 该操作应该是原子的,后续优化
	item := dto.Applied{
		UserId: friend.FriendUserId,
		Status: dto.Agree,
	}
	if err := document.UpdateAppliedStatus(friend.SelfUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	item = dto.Applied{
		UserId: friend.SelfUserId,
		Status: dto.OtherAgree,
	}
	if err := document.UpdateAppliedStatus(friend.FriendUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	// add friendId in friends  添加申请Id到好友列表
	if err := document.UpdateFriendsByUserId(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}

	// add self id in application user 添加自己到申请人的好友列表
	if err := document.UpdateFriendsByUserId(friend.FriendUserId, friend.SelfUserId); err != nil {
		return err
	}

	return nil
}

// NotPassFriend 未通过申请处理  主体是被申请人
func NotPassFriend(friend *friend_dto.ReqFriendAdd) error {
	// XXX 该操作应该是原子的,后续优化
	item := dto.Applied{
		UserId: friend.FriendUserId,
		Status: dto.NoPass,
	}
	if err := document.UpdateAppliedStatus(friend.SelfUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	item = dto.Applied{
		UserId: friend.SelfUserId,
		Status: dto.OtherNoPass,
	}
	if err := document.UpdateAppliedStatus(friend.FriendUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	return nil
}

// DeleteFriend  刪除好友,双向删除
func DeleteFriend(friend *friend_dto.ReqFriendAdd) error {
	// XXX \该操作应该是原子的,后续优化
	// I delete friend
	if err := document.DeleteFriend(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}
	// friend delete me
	if err := document.DeleteFriend(friend.FriendUserId, friend.SelfUserId); err != nil {
		return err
	}

	return nil
}

// FindSugFriends  推荐好友, 1.自己国家 2.积分大于自己  3.每次5个,不足补充其它国家的
func FindSugFriends(req *friend_dto.ReqRecommend) ([]*friend_dto.RespFriendRecommend, error) {

	const recommendNumber = 5

	// 查询自身
	user, err := document.SelectUser(int(req.UserId))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}

	// 查询符合条件的同国家玩家
	users, err := document.SelectFriendByCountryAndPoints(user.Country, user.Points, recommendNumber, false)
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}

	// 数量不足补充
	if recommendNumber-len(users) > 0 {
		supplementUsers, err := document.SelectFriendByCountryAndPoints(user.Country, user.Points, int64(recommendNumber-len(users)), true)
		if err != nil {
			return nil, errors.Wrap(err, response.MsgFailed)
		}
		users = append(users, supplementUsers...)
	}

	return users, nil
}

func ExistsInFriends(friend *friend_dto.ReqFriendAdd) error {
	user, err := document.SelectUser(int(friend.SelfUserId))
	if err != nil {
		return err
	}
	for _, v := range user.Friends {
		if v == friend.FriendUserId {
			return errors.New("already exists in friends list")
		}
	}
	return nil
}
