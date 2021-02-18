package main

import (
	"Notification/Backend/controller"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

func main() {
	//Load Cred mongodb
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Desktop Notification pop differnt for differ system
	if runtime.GOOS == "windows" {
		executeWindow()
		//process.Worker()
	} else {
		flag.Parse()
		log.SetFlags(0)
		http.HandleFunc("/echo", controller.Echo)
		http.HandleFunc("/", controller.HomePage)
		fmt.Println("Listening")
		log.Fatal(http.ListenAndServe(*addr, nil))
		//service.GetDetailForUser("milind")
		//userLogin()
	}
}

//executeWindow is used to send notification in windows machine using beep lib
func executeWindow() {
	err := beeep.Notify("Title", "Message body", "assets/information.png")
	if err != nil {
		panic(err)
	}
}

//go build -o main.sh
//chmod +x filename.sh
//./filename.sh
// https://askubuntu.com/questions/229589/how-to-make-a-file-e-g-a-sh-script-executable-so-it-can-be-run-from-a-termi

// func execute() {
// 	_, err := exec.Command("notify-send", "title", "body").Output()
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 	}
// 	fmt.Println("Command Successfully Executed")
// 	//output := string(out[:])
// 	//fmt.Println(output)
// }

// func userLogin() {
// 	user, err := user.Current()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Hello %s, Welcome to Desktop Notifier \n", user.Username)
// 	login.Start(os.Stdin, os.Stdout)

// }
