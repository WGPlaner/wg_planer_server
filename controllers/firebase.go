package controllers

import (
	"context"
	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/api/option"
	"os"
	"log"
)

var fireBaseApp *firebase.App
var fireBaseAuth *firebase.Auth

func InitialiseFirebaseConnection() {
	var err error

	if _, err := os.Stat("./config/serviceAccountKey.json"); os.IsNotExist(err) {
		log.Fatal("File is missing: config/serviceAccountKey.json") // exit program
	}

	fireBaseApp, err = firebase.InitializeApp(context.Background(), firebase.AppOptions{
		ProjectID: "wgplaner-se",
	}, option.WithCredentialsFile("./config/serviceAccountKey.json"))

	fireBaseAuth = fireBaseApp.Auth()

	if err != nil {
		panic(err)
	}
}
