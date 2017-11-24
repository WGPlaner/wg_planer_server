package setting

import (
	"context"
	"log"
	"path"

	"github.com/acoshift/go-firebase-admin"
	"github.com/wgplaner/wg_planer_server/modules/base"
	"google.golang.org/api/option"
)

var FireBaseApp *firebase.App

func CreateFirebaseConnection() *firebase.App {
	var (
		fireBaseApp *firebase.App
		err         error
	)

	keyfilePath := path.Join(AppWorkPath, "config/serviceAccountKey.json")
	base.FileMustExist(keyfilePath)

	fireBaseApp, err = firebase.InitializeApp(context.Background(), firebase.AppOptions{
		ProjectID: AppConfig.Auth.FirebaseProjectId,
		APIKey:    AppConfig.Auth.FirebaseServerKey,
	}, option.WithCredentialsFile(keyfilePath))

	if err != nil {
		log.Fatalln("[Firebase] Creation using key failed")
		return nil
	}

	return fireBaseApp
}
