package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserData struct {
	Email string
	Name  string
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	userData := &UserData{Email: "toto@gmail.com", Name: "Frank"}
	userDataStr, err := json.Marshal(userData)
	if err != nil {
		fmt.Fprintf(w, "{Error:%s}", err)
		return
	}
	fmt.Fprintf(w, string(userDataStr))
}

func SetupRoutes() {
	http.HandleFunc("/getUserData", getUserData)
}

func main() {
	portNumber := "9000"
	SetupRoutes()
	fmt.Println("Server listening on port ", portNumber)
	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}
