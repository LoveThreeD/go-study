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

/*
	从mongo中获取某位用户的全部数据
*/
func SelectUser(userId int) (c *dto.UserBaseData, err error) {
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

func SelectUserByNickname(reqParams *friend_dto.ReqFriendSearch) ([]*friend_dto.RespFriendRecommend, error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		mongo_key.BaseNickName: reqParams.NickName,
	}
	// sort
	sort := bson.M{mongo_key.BaseIsOnline: -1, mongo_key.BaseOfflineTime: -1}

	c := []*friend_dto.RespFriendRecommend{}
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if err = cursor.All(context.TODO(), &c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return c, nil
}

func AddPoints(userId int, points int) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
	}
	update := bson.M{
		"$inc": bson.M{
			mongo_key.BasePoints: points,
		},
	}
	if _, err = collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

// UpdateElementByUserId 更新该文档中的任一值
// 要是传入bson.M{}的话,在server层写db层的一些代码.不知道如何取舍，或者使用...可变参数
func UpdateElementByUserId(itemName string, userId int64, v interface{}) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
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
		return errors.New(response.MsgMongoUpdateUserError)
	}
	return nil
}

// 要是传入bson.M{}的话,在server层写db层的一些代码.不知道如何取舍，或者使用...可变参数
func UpdateTwoElementByUserId(userId int64, itemName1, itemName2 string, v1, v2 interface{}) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		mongo_key.BaseUserId: userId,
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
		return errors.New(response.MsgMongoUpdateUserError)
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
		return errors.New(response.MsgMongoUpdateUserError)
	}
	return nil
}

// SelectFriendByCountryAndPoints 好友推荐
func SelectFriendByCountryAndPoints(country string, integral int, limit int64, diffCountry bool) ([]*friend_dto.RespFriendRecommend, error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		mongo_key.BaseCountry: country,
		mongo_key.BasePoints: bson.M{
			"$gt": integral,
		},
	}

	// diffCountry is true. add Recommend friend
	if diffCountry {
		filter = bson.M{
			mongo_key.BaseCountry: bson.M{
				"$ne": country,
			},
			mongo_key.BasePoints: bson.M{
				"$gt": integral,
			},
		}
	}

	c := []*friend_dto.RespFriendRecommend{}
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if err = cursor.All(context.TODO(), &c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return c, nil
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
		"$addToSet": bson.M{
			mongo_key.BaseApplied: item,
		},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if result.ModifiedCount < 1 {
		return errors.New(response.MsgMongoUpdateUserError)
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
		return errors.New(response.MsgMongoUpdateUserError)
	}
	return nil
}

// return mongo collection connection
func getUserDocumentConnect() (c *mongo.Collection, err error) {
	return mongo_db.GetDocumentConnect("user", "base")
}
