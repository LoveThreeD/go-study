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
