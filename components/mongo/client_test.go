package mongo

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func NewMongoTest() (*Mongo, error) {
	hosts := "10.136.158.10:27000,10.136.158.10:27001,10.136.158.10:27002"
	replset := "rs_image"
	setting := Setting{
		Hosts:   hosts,
		ReplSet: replset,
	}
	Mongo, err := NewMongo(setting)
	if err != nil {
		return nil, err
	}
	return Mongo, nil
}

func TestNewMongo(t *testing.T) {
	hosts := "10.136.158.10:27000,10.136.158.10:27001,10.136.158.10:27002"
	replset := "rs_image"

	setting := Setting{
		Hosts:   hosts,
		ReplSet: replset,
	}
	Mongo, err := NewMongo(setting)
	if err != nil {
		t.Error("new mongo init error:", err)
		return
	}

	t.Log("init mongo success!")

	c := Mongo.Client.Database("test").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := c.Find(ctx, bson.D{})
	if err != nil {
		t.Error("find collection error:", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			t.Error("cur error")
			return
		}
		t.Log("result:", result)
	}

}

func TestOnlineMongo(t *testing.T) {
	hosts := "rs_image_status_01.int.yidian-inc.com:27017,rs_image_status_02.int.yidian-inc.com:27017,rs_image_status_03.int.yidian-inc.com:27017"
	replset := "rs_image_status"

	setting := Setting{
		Hosts:    hosts,
		ReplSet:  replset,
		Username: "rs_image_status_rw",
		Password: "w#JR9tQOVud$jh43V7v#B9Bd#kD4w%",
		//authSource: "admin",
	}
	Mongo, err := NewMongo(setting)
	if err != nil {
		t.Error("new mongo init error:", err)
		return
	}

	t.Log("init mongo success!")

	c := Mongo.Client.Database("test").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	res, err := c.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	id := res.InsertedID
	t.Log(id)

	defer cancel()
	cur, err := c.Find(ctx, bson.D{})
	if err != nil {
		t.Error("find collection error:", err)
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			t.Error("cur error")
			return
		}
		t.Log("result:", result)
	}

}
