package wgplaner

import (
	"context"
	"log"
	"path"

	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/api/option"
)

var FireBaseApp *firebase.App

func CreateFirebaseConnection() *firebase.App {
	var (
		fireBaseApp *firebase.App
		err         error
	)

	keyfilePath := path.Join(AppWorkPath, "config/serviceAccountKey.json")
	FileMustExist(keyfilePath)

	fireBaseApp, err = firebase.InitializeApp(context.Background(), firebase.AppOptions{
		ProjectID: "wgplaner-se",
	}, option.WithCredentialsFile(keyfilePath))

	if err != nil {
		log.Fatalln("[Firebase] Creation using key failed")
		return nil
	}

	return fireBaseApp
}
