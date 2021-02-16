package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gorilla/websocket"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ConnectTODatabase is used initalize connection to database
func ConnectTODatabase() *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("dbUrl")))
	if err != nil {
		panic(err)
	}
	//defer client.Disconnect(context.TODO())
	database := client.Database(os.Getenv("dbName"))
	return database
}

//MonitorStream is used to monitor collection
func MonitorStream(database *mongo.Database, c chan<- string, objectid interface{}) {
	collection := database.Collection(os.Getenv("dbCollection"))
	pipeline := bson.D{
		{
			Key: "$match", Value: bson.D{
				{Key: "fullDocument._id", Value: objectid}, //_id actual
			},
		},
	}
	collectionStream, errStream := collection.Watch(context.TODO(), mongo.Pipeline{pipeline})
	if errStream != nil {
		panic(errStream)
	}
	defer collectionStream.Close(context.TODO())
	for collectionStream.Next(context.TODO()) {
		var result map[string]interface{}
		if err := collectionStream.Decode(&result); err != nil {
			panic(err)

		}
		d := result["fullDocument"]
		dbByte, _ := json.Marshal(d)
		var fullDocument map[string]string
		_ = json.Unmarshal(dbByte, &fullDocument)
		fmt.Printf("%v\n", fullDocument["data"])
		//give call to the socket
		c <- fullDocument["data"]
	}
}

//Reader is used for monitoring chanages
func Reader(conn *websocket.Conn, objid interface{}) {
	c := make(chan string)
	db := ConnectTODatabase()
	//Get logined user _id and monitor that
	go MonitorStream(db, c, objid)
	for dchan := range c {
		b := []byte(dchan)
		if err := conn.WriteMessage(1, b); err != nil {
			log.Println(err)
			return
		}
	}
}

//Writer is used for writing to the frontend
func Writer(conn *websocket.Conn) {
	for {
		fmt.Println("Sending")
		messageType, r, err := conn.NextReader()
		if err != nil {
			fmt.Println(err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			fmt.Println(err)
			return
		}
		if err := w.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}
}

//ListenFormFrontEnd used to get value from frontend
func ListenFormFrontEnd(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error Listening", err)
			return
		}
		fmt.Println(">>>:", string(msg))

	}
}
