package dbhandler

import (
	"context"
	"core/utils"
	"fmt"
	"log"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func Connect() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("%s+srv://%s:%s@cluster0.qxgzt.%s.net/test", config.DbDriver, config.DBUsername, config.DBPassword, config.DbDriver)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	Db = client.Database(config.DBName)
}
