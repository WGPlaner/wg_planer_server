package wgplaner

import "testing"

func TestValidateDriverConfig(t *testing.T) {
	invalidConfig := databaseConfig{}
	if errs := ValidateDriverConfig(invalidConfig); !errs.HasErrors() {
		t.Error("Empty config should be an error!")
	}

	invalidConfig.Driver = "INVALID"
	if errs := ValidateDriverConfig(invalidConfig); len(errs.errors) != 1 {
		t.Error("Validation should exist after invalid driver check!")
	}

}
