package wgplaner

import (
	"context"
	"log"
	"os"

	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/api/option"
)

var FireBaseApp *firebase.App

func CreateFirebaseConnection() *firebase.App {

	var fireBaseApp *firebase.App
	var err error

	if _, err := os.Stat("./config/serviceAccountKey.json"); os.IsNotExist(err) {
		log.Fatal("File is missing: config/serviceAccountKey.json") // exit program
	}

	fireBaseApp, err = firebase.InitializeApp(context.Background(), firebase.AppOptions{
		ProjectID: "wgplaner-se",
	}, option.WithCredentialsFile("./config/serviceAccountKey.json"))

	if err != nil {
		log.Fatalln("[Firebase] Creation using key failed")
		return nil
	}

	return fireBaseApp

}
