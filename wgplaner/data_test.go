package wgplaner

import "testing"

func TestValidateDataConfig(t *testing.T) {
	invalidConfig := dataConfig{}
	if errs := ValidateDataConfig(invalidConfig); !errs.HasErrors() {
		t.Error("Empty data config should be an error!")
	}

	invalidConfig.UserImageDir = "does_not_exist"
	if errs := ValidateDataConfig(invalidConfig); len(errs.errors) != 1 {
		t.Error("Validation should exit after invalid data check!")
	}

	invalidConfig.UserImageDir = "../config/config.example.toml"
	if errs := ValidateDataConfig(invalidConfig); len(errs.errors) != 1 {
		t.Error("File as 'user_image_dir' should throw an error")
	}
}
