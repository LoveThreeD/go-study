package document

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sTest/entity/dto"
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

// return mongo collection connection
func getUserDocumentConnect() (c *mongo.Collection, err error) {
	return mongo_db.GetDocumentConnect("user", "base")
}
