package main

import (
	"fmt"
	"github.com/jiaying2001/agent/auth"
	"github.com/jiaying2001/agent/launcher"
	"github.com/jiaying2001/agent/myzk"
	"github.com/jiaying2001/agent/store"
	"os"
)

func login() {
	var username, password string
	// Prompt the user to enter username and password
	fmt.Print("Enter username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		return
	}
	fmt.Print("Enter password: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		return
	}
	if auth.Login(username, password) {
		fmt.Println("Successfully logged in")
	} else {
		fmt.Println("no ma mei si")
		os.Exit(1)
	}
}

func main() {
	login()
	myzk.LoadIdsNodes()
	myzk.ListenIdsNodes()
	myzk.Listen("/" + store.Pass.UserName) // Listen a ZK node
	launcher.L.Launch()
	select {}
}

//func main() {
//}
