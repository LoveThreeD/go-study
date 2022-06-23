package service

import (
	"github.com/pkg/errors"
	"sTest/entity/dto"
	"sTest/entity/friend_dto"
	"sTest/pkg/response"
	"sTest/repository/document"
)

func SearchUser(search *friend_dto.ReqFriendSearch) (c []dto.UserBaseData, err error) {
	users, err := document.SelectUserByNickName(search)
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}
	return users, nil
}

// AddApplicationList  添加到申请列表
func AddApplicationList(friend *friend_dto.ReqFriendAdd) (err error) {
	if err = document.UpdateApplicationListByUserId(friend.FriendUserId, friend.SelfUserId); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}
	return nil
}

// AddFriendList  添加到好友列表
func AddFriendList(friend *friend_dto.ReqFriendAdd) (err error) {
	// XXX \该操作应该是原子的,后续优化
	// delete application list friendId      删除申请列表中的申请Id
	if err = document.DeleteApplicationList(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}

	// add friendId in friends  添加申请Id到好友列表
	if err = document.UpdateFriendsByUserId(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}

	// add self id in application user 添加自己到申请人的好友列表
	if err = document.UpdateFriendsByUserId(friend.FriendUserId, friend.SelfUserId); err != nil {
		return err
	}

	return nil
}

// DeleteFriendList  刪除好友,双向删除
func DeleteFriendList(friend *friend_dto.ReqFriendAdd) (err error) {
	// XXX \该操作应该是原子的,后续优化
	// I delete friend
	if err = document.DeleteFriendList(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}
	// friend delete me
	if err = document.DeleteFriendList(friend.FriendUserId, friend.SelfUserId); err != nil {
		return err
	}

	return nil
}
