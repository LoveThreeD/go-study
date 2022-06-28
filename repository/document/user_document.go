package document

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sTest/entity/dto"
	"sTest/entity/friend_dto"
	"sTest/pkg/mongo_db"
	"sTest/pkg/response"
	"sTest/repository/document/mongo_key"
)

// CreateUser create one user
func CreateUser(item *dto.UserBaseData) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return err
	}
	if _, err := collection.InsertOne(context.TODO(), item); err != nil {
		return errors.Wrap(err, response.MsgMongoCreateUserError)
	}
	return
}

func SelectUserByUserId(userId int64) (c *dto.UserCache, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
	}
	item := &dto.UserBaseData{}
	if err = collection.FindOne(context.TODO(), filter).Decode(item); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	c = &dto.UserCache{}
	c.NickName = item.NickName
	c.AvatarUrl = item.AvatarURL
	c.IsOnline = 1
	return
}

func SelectUserByUserIdAll(userId int) (c *dto.UserBaseData, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
	}
	c = &dto.UserBaseData{}
	if err = collection.FindOne(context.TODO(), filter).Decode(c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

func SelectUserByNickName(reqParams *friend_dto.ReqFriendSearch) (c []friend_dto.RespFriendRecommend, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		mongo_key.BaseNickName: reqParams.NickName,
	}
	// sort
	sort := bson.D{{mongo_key.BaseIsOnline, -1}, {mongo_key.BaseOfflineTime, -1}}

	c = []friend_dto.RespFriendRecommend{}
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if err = cursor.All(context.TODO(), &c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

func IncrIntegral(userId int, integral int) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$inc": bson.M{
			"integral": integral,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

// XXX 通过userId更新数组可以封装

func UpdateApplicationListByUserId(userId int64, applicationId int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$push": bson.M{
			"applicationlist": applicationId,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return nil
}

func UpdateAlreadyAppliedListByUserId(userId int64, alreadyApplied int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$push": bson.M{
			"alreadyappliedlist": alreadyApplied,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return nil
}

func UpdateNoPassListByUserId(userId int64, noPass int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$push": bson.M{
			"nopass": noPass,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return nil
}

// UpdateAddElementByUserId 通过userId查找文档,添加文档中arrayName数组的元素
func UpdateAddElementByUserId(arrayName string, userId int64, noPass int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$push": bson.M{
			arrayName: noPass,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return nil
}

// DeleteElementByUserId 通过userId查找文档,删除文档中arrayName数组的元素
func DeleteElementByUserId(arrayName string, userId int64, friendId int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$pull": bson.M{
			arrayName: friendId,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

// UpdateElementByUserId 更新该文档中的任一值
func UpdateElementByUserId(itemName string, userId int64, v interface{}) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$set": bson.M{
			itemName: v,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

func UpdateTwoElementByUserId(userId int64, itemName1, itemName2 string, v1, v2 interface{}) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$set": bson.M{
			itemName1: v1,
			itemName2: v2,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

func UpdateFriendsByUserId(userId int64, friendId int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
	}
	update := bson.M{
		"$push": bson.M{
			mongo_key.BaseFriends: friendId,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return nil
}

// DeleteApplicationList 删除申请列表中的userId
func DeleteApplicationList(userId int64, friendId int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$pull": bson.M{
			"applicationlist": friendId,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

// DeleteFriendList 删除好友列表中的userId
func DeleteFriendList(userId int64, friendId int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
	}
	update := bson.M{
		"$pull": bson.M{
			mongo_key.BaseFriends: friendId,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

// SelectFriendByCountryAndIntegral 好友推荐
func SelectFriendByCountryAndIntegral(country string, integral int, limit int64, diffCountry bool) (c []friend_dto.RespFriendRecommend, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		mongo_key.BaseCountry: country,
		mongo_key.BaseIntegral: bson.M{
			"$gt": integral,
		},
	}

	// diffCountry is true. add Recommend friend
	if diffCountry {
		filter = bson.M{
			mongo_key.BaseCountry: bson.M{
				"$ne": country,
			},
			mongo_key.BaseIntegral: bson.M{
				"$gt": integral,
			},
		}
	}

	c = []friend_dto.RespFriendRecommend{}
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if err = cursor.All(context.TODO(), &c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

// AddApplied 添加申请 保证 列表中没有userId == userId 以及 status 为 申请0、被申请1
func AddApplied(userId int64, item *dto.Applied) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
	}
	update := bson.M{
		"$push": bson.M{
			mongo_key.BaseApplied: item,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

// UpdateAppliedStatus 修改申请状态
func UpdateAppliedStatus(userId int64, item *dto.Applied) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId:          userId,
		mongo_key.BaseUserIdInApplied: item.UserId,
	}
	update := bson.M{
		"$set": bson.M{
			mongo_key.BaseUpStatusInApplied: item.Status,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrap(errors.New(response.MsgMongoUpdateUserError), response.MsgMongoUpdateUserError)
	}
	return nil
}

// return mongo collection connection
func getUserDocumentConnect() (c *mongo.Collection, err error) {
	return mongo_db.GetDocumentConnect("user", "base")
}
