package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("error: ", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("error : ", err)
	}
	fmt.Println("Connected to MongoDb	 ")
	return client
}
