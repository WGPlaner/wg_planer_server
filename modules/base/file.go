package base

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func GetMimeType(data []byte) string {
	return http.DetectContentType(data)
}

func IsFileJPG(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "image/jpeg")
}

// IsFileImage detects if data is an image format
func IsFileImage(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "image/")
}

func FileMustExist(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf(`File is missing: "%s"`, filePath)

	} else if os.IsPermission(err) {
		log.Fatalf(`Permissions error for "%s"`, filePath)

	} else if err != nil {
		log.Fatalf(`Unknown error reading "%s"`, filePath)
	}
}
