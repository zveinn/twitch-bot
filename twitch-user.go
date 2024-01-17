package main

import (
	"context"
	"log"
	"time"

	tirc "github.com/gempir/go-twitch-irc/v4"
	"github.com/zveinn/twitch-bot/mongowrapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Points      int            `json:"Points" bson:"Points"`
	ID          string         `json:"ID" bson:"ID"`
	Name        string         `json:"Name" bson:"Name"`
	DisplayName string         `json:"DisplayName" bson:"DisplayName"`
	Color       string         `json:"Color" bson:"Color"`
	Badges      map[string]int `json:"Badges" bson:"Badges"`
}

type UserMSG struct {
	Raw            string            `bson:"Raw"`
	Type           tirc.MessageType  `bson:"Type"`
	RawType        string            `bson:"RawType"`
	Tags           map[string]string `bson:"Tags"`
	Message        string            `bson:"Message"`
	Channel        string            `bson:"Channel"`
	RoomID         string            `bson:"RoomID"`
	ID             string            `bson:"ID"`
	Time           time.Time         `bson:"Time"`
	Emotes         []*tirc.Emote     `bson:"Emotes"`
	Bits           int               `bson:"Bits"`
	Action         bool              `bson:"Action"`
	FirstMessage   bool              `bson:"FirstMessage"`
	Reply          *tirc.Reply       `bson:"Reply"`
	CustomRewardID string            `bson:"CustomRewardID"`
}

func GetTop10() (userList []*User) {
	opts := options.Find().SetSort(bson.D{{"Points", -1}}).SetLimit(11)
	ctx := context.Background()

	cursor, err := mongowrapper.UserCollection.Find(
		ctx,
		bson.D{},
		opts,
	)
	if err != nil {
		log.Println("Unable to decode top10", err)
		return
	}

	userList = make([]*User, 0)
	err = cursor.All(ctx, &userList)
	if err != nil {
		log.Println("Unable to decode top10", err)
		return
	}

	return
}

func IncrementUserPoints(user *User, points int) (err error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"uid": user.ID}
	ctx := context.Background()

	err = mongowrapper.UserCollection.FindOneAndUpdate(
		ctx,
		filter,
		bson.D{
			{"$inc", bson.D{{"Points", points}}},
		},
		opts,
	).Err()

	if err != nil {
		log.Println("ERROR INCREMENTING USER STATS", err)
	}

	return
}

func FindUserMessagesFromMatch(user string, match string) (msgList []*tirc.PrivateMessage, err error) {
	opts := options.Find()
	filter := bson.D{
		{"message", primitive.Regex{Pattern: match, Options: ""}},
		{"user.name", user},
	}
	ctx := context.Background()

	cursor, err := mongowrapper.MSGCollection.Find(
		ctx,
		filter,
		opts,
	)
	if err != nil {
		log.Println("Error getting quote", err)
		return
	}

	msgList = make([]*tirc.PrivateMessage, 0)
	err = cursor.All(ctx, &msgList)
	if err != nil {
		log.Println("Error parsing quote:", err)
		return
	}

	return
}

func FindOrUpsertUser(user *tirc.User) (U *User, err error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"uid": user.ID}

	ctx := context.Background()

	U = new(User)
	err = mongowrapper.UserCollection.FindOneAndUpdate(
		ctx,
		filter,
		bson.D{
			{"$set", bson.D{{"lastSeen", time.Now().UnixNano()}}},
			{"$set", bson.D{{"ID", user.ID}}},
			{"$set", bson.D{{"Name", user.Name}}},
			{"$set", bson.D{{"DisplayName", user.DisplayName}}},
			{"$set", bson.D{{"Color", user.Color}}},
			{"$set", bson.D{{"Badges", user.Badges}}},
		},
		opts,
	).Decode(&U)
	if err != nil {
		log.Println("ERROR FINDING USER:", err)
		return
	}

	return
}

func SaveMessage(msg *tirc.PrivateMessage) (err error) {
	ctx := context.Background()
	_, err = mongowrapper.MSGCollection.InsertOne(ctx, msg, options.InsertOne().SetBypassDocumentValidation(true))
	if err != nil {
		log.Println("ERROR INSERTING MSG")
		return
	}
	return
}
