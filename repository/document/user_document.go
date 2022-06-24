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

func SelectUserByUserId(userId int) (c *dto.UserCache, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"userid":   userId,
		"basedata": 1,
	}
	item := &dto.UserBaseData{}
	if err = collection.FindOne(context.TODO(), filter).Decode(item); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	c = &dto.UserCache{}
	c.NickName = item.BaseData.NickName
	c.AvatarUrl = item.BaseData.AvatarURL
	c.IsOnline = 1
	return
}

func SelectUserByUserIdAll(userId int) (c *dto.UserBaseData, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"userid": userId,
	}
	c = &dto.UserBaseData{}
	if err = collection.FindOne(context.TODO(), filter).Decode(c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

func SelectUserByNickName(reqParams *friend_dto.ReqFriendSearch) (c []dto.UserBaseData, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"basedata.nickname": reqParams.NickName,
	}
	c = []dto.UserBaseData{}
	cursor, err := collection.Find(context.TODO(), filter)
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

func UpdateFriendsByUserId(userId int64, friendId int64) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return
	}
	filter := bson.M{
		"userid": userId,
	}
	update := bson.M{
		"$push": bson.M{
			"friends": friendId,
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
		"userid": userId,
	}
	update := bson.M{
		"$pull": bson.M{
			"friends": friendId,
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

func SelectFriendByCountryAndIntegral(country string, integral int, limit int64, diffCountry bool) (c []dto.UserBaseData, err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"country": country,
		"integral": bson.M{
			"$gt": integral,
		},
	}

	// diffCountry is true. add Recommend friend
	if diffCountry {
		filter = bson.M{
			"country": bson.M{
				"$ne": country,
			},
			"integral": bson.M{
				"$gt": integral,
			},
		}
	}

	c = []dto.UserBaseData{}
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	if err = cursor.All(context.TODO(), &c); err != nil {
		return nil, errors.Wrap(err, response.MsgMongoSelectUserError)
	}
	return
}

// return mongo collection connection
func getUserDocumentConnect() (c *mongo.Collection, err error) {
	return mongo_db.GetDocumentConnect("user", "base")
}
