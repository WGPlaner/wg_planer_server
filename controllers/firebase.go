package controllers

import (
	"context"
	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/api/option"
)

var fireBaseApp *firebase.App
var fireBaseAuth *firebase.Auth

func InitialiseFirebaseConnection() {
	var err error

	fireBaseApp, err = firebase.InitializeApp(context.Background(), firebase.AppOptions{
		ProjectID: "wgplaner-se",
	}, option.WithCredentialsFile("serviceAccountKey.json"))

	fireBaseAuth = fireBaseApp.Auth()

	if err != nil {
		panic(err)
	}
}
