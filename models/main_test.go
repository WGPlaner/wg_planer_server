package models

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" // for the test engine
)

func TestMain(m *testing.M) {
	if err := CreateTestEngine("fixtures/"); err != nil {
		fmt.Printf("Error creating test engine: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
