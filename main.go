package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Checking github ssh again
	mongoConn := options.Client().ApplyURI("mongodb://127.0.0.1:27017/")
	ctx := context.Background()
	client, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	defer client.Disconnect(ctx)
	// проверка связи с БД
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	dbs, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DATABASES: ", dbs)

	mogodb := client.Database("mogodb")
	languages := mogodb.Collection("languages")

	result, err := languages.InsertOne(ctx, bson.D{{Key: "ID", Value: 1}, {Key: "name", Value: "Golang"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result:", result)
	defer languages.Drop(ctx)

	results, err := languages.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "ID", Value: 2},
			{Key: "name", Value: "Java"},
		},
		bson.D{
			{Key: "ID", Value: 3},
			{Key: "name", Value: "JavaScript"},
		},
		bson.D{
			{Key: "ID", Value: 4},
			{Key: "name", Value: "Python"},
		},
		bson.D{
			{Key: "ID", Value: 5},
			{Key: "name", Value: "C++"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Results:", results)

	cursor, err := languages.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var lang []bson.M
	if err := cursor.All(ctx, &lang); err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	print(lang)
	//fmt.Println("Cursor:", lang)
	for cursor.Next(ctx) {
		var l bson.M
		if err := cursor.Decode(&l); err != nil {
			log.Fatal(err)
		}
		lang = append(lang, l)
	}
	fmt.Println()
	print(lang)

	var sprache bson.M
	if err := languages.FindOne(ctx, bson.M{}).Decode(&sprache); err != nil {
		log.Fatal(err)
	}
	fmt.Println("language:", sprache)

	filter := bson.M{"name": "Golang"}
	filterCursor, err := languages.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var lan []bson.M
	if err := filterCursor.All(ctx, &lan); err != nil {
		log.Fatal(err)
	}
	fmt.Println("filtered Golang:")
	print(lan)
	fmt.Println("###################################")
	gonews := client.Database("gonews")
	fmt.Println("gonews DATABASE ")
	gonewsColl := gonews.Collection("news")
	fmt.Println("gonews Collection news ????")
	res, _ := gonews.ListCollectionNames(ctx, bson.D{{"options.capped", true}})

	fmt.Println("gonews: ", res)
	cur, err := gonewsColl.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("gonewsColl ERROR ")
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		fmt.Println("next")
	}

}
func print(a []bson.M) {
	for _, val := range a {
		fmt.Println(val)
	}
}
