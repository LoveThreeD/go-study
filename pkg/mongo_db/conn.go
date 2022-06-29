package mongo_db

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sTest/pkg/response"
	"sTest/pkg/viper"
)

// ClientOpts mongoClient 连接客户端参数
var ClientOpts *options.ClientOptions
var Client *mongo.Client

func init() {
	var err error
	c := viper.Conf.Mongo
	url := fmt.Sprintf("mongodb://%s:%d", c.Address, c.Port)
	ClientOpts = options.Client().ApplyURI(url)
	Client, err = mongo.Connect(context.TODO(), ClientOpts)
	if err != nil {
		logger.Fatal(err)
	}
}

func GetDocumentConnect(dbName, collName string) (c *mongo.Collection, err error) {
	db := Client.Database(dbName)
	if db == nil {
		return nil, errors.New(response.MsgMongoDbConnectionError)
	}
	c = db.Collection(collName)
	if c == nil {
		return nil, errors.New(response.MsgMongoCollConnectionError)
	}
	return c, nil
}
