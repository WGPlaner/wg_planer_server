package wgplaner

import (
	"os"
	"path"

	"github.com/wgplaner/wg_planer_server/gen/models"

	"github.com/go-openapi/swag"
)

func ValidateDataConfig(config dataConfig) ErrorList {
	errList := ErrorList{}

	if stat, err := os.Stat(path.Join(AppWorkPath, config.UserImageDir)); err != nil {
		if os.IsNotExist(err) {
			errList.Add("[Config][Data] 'user_image_dir' does not exist!")
		} else if os.IsPermission(err) {
			errList.Add("[Config][Data] Permission denied for 'user_image_dir'!")
		} else {
			errList.Add("[Config][Data] Unknown error with 'user_image_dir'! " + err.Error())
		}

	} else if !stat.IsDir() {
		errList.Add("[Config][Data] 'user_image_dir' is not a directory!")
	}

	return errList
}

const (
	PROFILE_IMAGE_FILE_NAME = "profile_image.jpg"
)

func GetUserProfileImageFilePath(user *models.User) string {
	return path.Join(
		AppConfig.Data.UserImageDir,
		swag.StringValue(user.UID),
		PROFILE_IMAGE_FILE_NAME,
	)
}

func GetUserProfileImage(user *models.User) (*os.File, error) {
	return os.Open(GetUserProfileImageFilePath(user))
}

func GetUserProfileImageDefault() (*os.File, error) {
	return os.Open(AppConfig.Data.UserImageDefault)
}
