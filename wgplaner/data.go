package wgplaner

import (
	"os"
	"path"

	"github.com/wgplaner/wg_planer_server/gen/models"

	"github.com/go-openapi/swag"
)

func ValidateDataConfig(config dataConfig) ErrorList {
	errList := ErrorList{}

	if config.UserImageDir == "" {
		errList.Add("[Config][Data] 'user_image_dir' must not be empty!")

	} else if stat, err := os.Stat(path.Join(AppWorkPath, config.UserImageDir)); err != nil {
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
	PROFILE_IMAGE_FILE_NAME       = "profile_image.jpg"
	GROUP_PROFILE_IMAGE_FILE_NAME = "group_image.jpg"
)

func GetUserProfileImageFilePath(user *models.User) string {
	return path.Join(
		AppWorkPath,
		AppConfig.Data.UserImageDir,
		swag.StringValue(user.UID),
		PROFILE_IMAGE_FILE_NAME,
	)
}

func GetUserProfileImage(user *models.User) (*os.File, error) {
	return os.Open(GetUserProfileImageFilePath(user))
}

func GetUserProfileImageDefault() (*os.File, error) {
	return os.Open(path.Join(
		AppWorkPath,
		AppConfig.Data.UserImageDefault,
	))
}

func GetGroupProfileImageFilePath(group *models.Group) string {
	return path.Join(
		AppWorkPath,
		AppConfig.Data.GroupImageDir,
		string(group.UID),
		GROUP_PROFILE_IMAGE_FILE_NAME,
	)
}

func GetGroupProfileImage(group *models.Group) (*os.File, error) {
	return os.Open(GetGroupProfileImageFilePath(group))
}

func GetGroupProfileImageDefault() (*os.File, error) {
	return os.Open(path.Join(
		AppWorkPath,
		AppConfig.Data.UserImageDefault,
	))
}
