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

func SendPushUpdateToUsers(users []*models.User, t PushUpdateType, data []string) error {
	if setting.AppConfig.Auth.IgnoreFirebase {
		return nil
	}

	var receiverIDs []string
	for _, u := range users {
		receiverIDs = append(receiverIDs, u.FirebaseInstanceID)
	}

	resp, err := setting.FireBaseApp.FCM().SendToDevices(context.Background(), receiverIDs, firebase.Message{
		Data: PushUpdateData{
			Type:    t,
			Updated: data,
		},
	})

	if err != nil {
		fireLog.Debug(`Error sending firebase update.`)
		return err
	}

	fireLog.Debug(resp)
	return nil
}

func SendPushUpdateToUserIDs(receiverIDs []string, t PushUpdateType, data []string) error {
	fireLog.Debug(`Send a firebase update data message to users (ids)`)

	if setting.AppConfig.Auth.IgnoreFirebase {
		return nil
	}

	users := make([]*models.User, 0, 10)

	for _, id := range receiverIDs {
		u, err := models.GetUserByUID(id)
		if err != nil {
			return err
		}
		users = append(users, u)
	}

	return SendPushUpdateToUsers(users, t, data)
}
