package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"game_genXls/lib"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	savePath  = flag.String("savePath", "", "Path to save the makefile")
	readPath  = flag.String("readPath", "", "The path of reading Excel")
	genStruct = flag.Int64("genStruct", 0, "gen struct | gen json")
	mongoUri  = flag.String("mongoUri", "", "mongo uri to save config data")
	mongoDb   = flag.String("mongoDb", "", "mongo db name to save config data")
)

func main() {
	flag.Parse()

	if *savePath == "" || *readPath == "" {
		fmt.Println("savePath or readPath  is nil")
		return
	}
	gt := lib.Generate{}
	if *genStruct == 1 {
		err := gt.GenStruct(*readPath, *savePath)
		if err != nil {
			fmt.Printf("generate struct err:%v\n", err)
			return
		}
	} else {
		if *mongoUri == "" || *mongoDb == "" {
			fmt.Println("mongoUri or mongoDb is nil")
		}
		clientOptions := options.Client().ApplyURI(*mongoUri)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = client.Database(*mongoDb).Drop(context.Background())
		if err != nil {
			log.Fatal(err)
			return
		}
		p := *savePath + "/json"
		err = gt.WriteMongo(client, *readPath, p, *mongoDb)
		if err != nil {
			fmt.Printf("generate json err:%v\n", err)
			return
		}
	}
}
