package dbutil

import (
	"Penggajian/pkg/util"
	"errors"
	"os"
	"strings"
)

func GetURI() (string, error) {
	var err error = nil

	mongoProtocol := os.Getenv("MONGODB_PROTOCOL")
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPass := os.Getenv("MONGODB_PASS")
	mongoUrl := os.Getenv("MONGODB_URL")
	mongoProps := os.Getenv("MONGODB_PROPERTIES")

	if util.IsEmpty(mongoProtocol) ||
		util.IsEmpty(mongoUser) ||
		util.IsEmpty(mongoPass) ||
		util.IsEmpty(mongoUrl) ||
		util.IsEmpty(mongoProps) {
		return "", errors.New("environment variable is not satisfied")
	}
	// Retrieve environment variable
	mongoUri := mongoProtocol + "://" + mongoUser + ":" + mongoPass + "@" + mongoUrl + "/" + mongoProps

	return mongoUri, err
}

func CheckURIFormat(uri_ string) bool {
	return strings.Index(uri_, "{USER}") != -1 && strings.Index(uri_, "{PASS}") != -1
}
