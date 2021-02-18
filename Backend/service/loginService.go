package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//GetDetailForUser is used to get _id based on uname and password
func GetDetailForUser(username string, password string) interface{} {
	db := ConnectTODatabase()
	loginHandler(db, username, password)
	result := loginHandler(db, username, password)
	if result != nil {
		fmt.Println("Login succeed")
		//service.MonitorStream(db, result["_id"])
		return result["_id"]
	}
	return nil
}

func loginHandler(database *mongo.Database, username string, password string) primitive.M {
	var result primitive.M
	collection := database.Collection("streamCollection")
	err := collection.FindOne(context.TODO(),
		bson.D{{Key: "name", Value: username}, {Key: "password", Value: password}}).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}

// -----Accepting Username and password through CMD------
// func Start(in io.Reader, out io.Writer) {
// 	//TODO CLI Accept login and pass and start process
// 	reader := bufio.NewReader(in)
// 	fmt.Print("Username: ")
// 	uname, _ := reader.ReadString('\n')
// 	fmt.Print("Password: ")
// 	upass, _ := reader.ReadString('\n')
// 	uname = strings.Replace(uname, "\n", "", -1)
// 	upass = strings.Replace(upass, "\n", "", -1)
// 	getDetailForUser(uname, upass)
// }
//--------------------------------------------------------
