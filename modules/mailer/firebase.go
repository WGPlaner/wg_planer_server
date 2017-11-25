package mailer

import (
	"context"

	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/setting"

	"github.com/acoshift/go-firebase-admin"
	"github.com/op/go-logging"
)

var fireLog = logging.MustGetLogger("Fire")

type (
	PushUpdateType string
	PushUpdateData struct {
		Type    PushUpdateType
		Updated []string
	}
)

const (
	PushUpdateGroup          = PushUpdateType("Group")
	PushUpdateGroupNewMember = PushUpdateType("GroupNewMember")
	PushUpdateUser           = PushUpdateType("User")
	PushUpdateShoppingList   = PushUpdateType("ShoppingList")
)

func SendPushUpdateToUsers(users []*models.User, t PushUpdateType, s []string) error {
	if setting.AppConfig.Auth.IgnoreFirebase {
		return nil
	}

	var ids []string
	for _, u := range users {
		ids = append(ids, u.FirebaseInstanceID)
	}

	return SendPushUpdateToUserIDs(ids, t, s)
}

func SendPushUpdateToUserIDs(reseiverIds []string, t PushUpdateType, data []string) error {
	fireLog.Debug(`Send a firebase update data message to users (ids)`)

	if setting.AppConfig.Auth.IgnoreFirebase {
		return nil
	}

	resp, err := setting.FireBaseApp.FCM().SendToDevices(context.Background(), reseiverIds, firebase.Message{
		Data: PushUpdateData{
			Type:    t,
			Updated: data,
		},
	})

	if err != nil {
		fireLog.Debug(`Error sending firebase update.`)
	}

	fireLog.Debug(resp)

	return err

}
