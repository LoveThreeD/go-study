package document

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"sTest/entity"
	"sTest/pkg/mongo_db"
	"sTest/pkg/response"
)

// CreateUser create one user
func CreateUser(item *entity.UserBaseData) (err error) {
	collection, err := getUserDocumentConnect()
	if err != nil {
		return err
	}
	if _, err := collection.InsertOne(context.TODO(), item); err != nil {
		return errors.Wrap(err, response.MsgMongoCreateUserError)
	}
	return
}

// return mongo collection connection
func getUserDocumentConnect() (c *mongo.Collection, err error) {
	return mongo_db.GetDocumentConnect("user", "base")
}
