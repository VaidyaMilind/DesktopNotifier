package controller

import (
	"Notification/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Init of Gorilla Websocket

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

//HomePage is used for displaying default home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

//Echo is used for forming actual connection with websocket and db
func Echo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uname, ok1 := query["username"]
	pass, ok2 := query["password"]

	//cross-origin check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if (!ok1 || len(uname[0]) < 1) || (!ok2 || len(pass[0]) < 1) {
		log.Println("Param are missing")
		return
	}

	//Getting _id of currently logged in user
	objectid := service.GetDetailForUser(uname[0], pass[0])
	if objectid == nil {
		log.Print("Failed while logging")
		return
	}
	fmt.Println("Id for given username and password is : ", objectid)

	go service.Reader(conn, objectid)
	go service.ListenFormFrontEnd(conn)
}
