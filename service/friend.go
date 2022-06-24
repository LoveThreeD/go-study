package service

import (
	"github.com/pkg/errors"
	"sTest/entity/dto"
	"sTest/entity/friend_dto"
	"sTest/pkg/response"
	"sTest/repository/document"
)

const (
	/*我不给通过的列表名*/
	_noPassName = "nopass"
	/*我申请后未通过的列表名*/
	_noPassRecordName = "nopassrecord"

	/*已申请列表(我申请别人的)*/
	_alreadyAppliedName = "alreadyappliedlist"
	/*被人申请我的申请列表*/
	_applicationListName = "applicationlist"
)

func SearchUser(search *friend_dto.ReqFriendSearch) (c []dto.UserBaseData, err error) {
	users, err := document.SelectUserByNickName(search)
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}
	return users, nil
}

// AddApplicationList  添加到申请列表 And 添加到已申请列表
func AddApplicationList(friend *friend_dto.ReqFriendAdd) (err error) {
	/*if err = document.UpdateApplicationListByUserId(friend.FriendUserId, friend.SelfUserId); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	if err = document.UpdateAlreadyAppliedListByUserId(friend.SelfUserId, friend.FriendUserId); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}*/

	item := dto.Applied{
		UserId: friend.FriendUserId,
		Status: dto.Apply,
	}
	if err = document.AddApplied(friend.SelfUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	item = dto.Applied{
		UserId: friend.SelfUserId,
		Status: dto.OtherApply,
	}
	if err = document.AddApplied(friend.FriendUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	return nil
}

// AddFriendList  添加到好友列表
func AddFriendList(friend *friend_dto.ReqFriendAdd) (err error) {
	// XXX 该操作应该是原子的,后续优化
	/*// delete application list friendId      删除已申请列表中的申请Id
	if err = document.DeleteElementByUserId(_alreadyAppliedName,friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}

	// add friendId in friends  添加申请Id到好友列表
	if err = document.UpdateFriendsByUserId(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}

	// add self id in application user 添加自己到申请人的好友列表
	if err = document.UpdateFriendsByUserId(friend.FriendUserId, friend.SelfUserId); err != nil {
		return err
	}*/

	item := dto.Applied{
		UserId: friend.FriendUserId,
		Status: dto.Agree,
	}
	if err = document.UpdateAppliedStatus(friend.SelfUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	item = dto.Applied{
		UserId: friend.SelfUserId,
		Status: dto.OtherAgree,
	}
	if err = document.UpdateAppliedStatus(friend.FriendUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
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

// NotPass 未通过申请处理  主体是被申请人
func NotPass(friend *friend_dto.ReqFriendAdd) (err error) {
	// XXX 该操作应该是原子的,后续优化

	/*// delete application list
	if err = document.DeleteApplicationList(friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}
	// add noPass list
	if err = document.UpdateAddElementByUserId(_noPassName,friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}

	// add noPassRecord list
	if err = document.UpdateAddElementByUserId(_noPassRecordName,friend.SelfUserId, friend.FriendUserId); err != nil {
		return err
	}*/

	item := dto.Applied{
		UserId: friend.FriendUserId,
		Status: dto.NoPass,
	}
	if err = document.UpdateAppliedStatus(friend.SelfUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}

	item = dto.Applied{
		UserId: friend.SelfUserId,
		Status: dto.OtherNoPass,
	}
	if err = document.UpdateAppliedStatus(friend.FriendUserId, &item); err != nil {
		return errors.Wrap(err, response.MsgFailed)
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

// GetRecommendFriends  推荐好友, 1.自己国家 2.积分大于自己  3.每次5个,不足补充其它国家的
func GetRecommendFriends(req *friend_dto.ReqRecommend) (f []dto.UserBaseData, err error) {

	const recommendNumber = 5

	// 查询自身
	user, err := document.SelectUserByUserIdAll(int(req.UserId))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}

	// 查询符合条件的同国家玩家
	users, err := document.SelectFriendByCountryAndIntegral(user.Country, user.Integral, recommendNumber, false)
	if err != nil {
		return nil, errors.Wrap(err, response.MsgFailed)
	}

	// 数量不足补充
	if recommendNumber-len(users) > 0 {
		supplementUsers, err := document.SelectFriendByCountryAndIntegral(user.Country, user.Integral, int64(recommendNumber-len(users)), true)
		if err != nil {
			return nil, errors.Wrap(err, response.MsgFailed)
		}
		users = append(users, supplementUsers...)
	}

	return users, nil
}
