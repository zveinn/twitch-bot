package mongowrapper

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client = MONGO{}

var (
	UserCollection *mongo.Collection
	MSGCollection  *mongo.Collection
)

type MONGO struct {
	Connection *mongo.Client
}

func InitCollections() {
	UserCollection = Client.Connection.Database("bot").Collection("users")
	MSGCollection = Client.Connection.Database("bot").Collection("msg")
}

func Connect(uri string) (err error) {
	var maxSize uint64 = 200
	var minSize uint64 = 20
	minHeartbeat := time.Duration(1 * time.Second)
	opt := options.Client()
	opt.MaxPoolSize = &maxSize
	opt.MinPoolSize = &minSize
	opt.HeartbeatInterval = &minHeartbeat
	Client.Connection, err = mongo.Connect(context.TODO(), opt.ApplyURI(uri))
	return err
}

func Disconnect() (err error) {
	err = Client.Connection.Disconnect(context.TODO())
	return err
}
