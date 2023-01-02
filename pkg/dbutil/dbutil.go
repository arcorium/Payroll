package dbutil

import (
	"errors"
	"os"
	"strings"
)

func GetURI() (string, error) {
	var err error = nil
	// Retrieve environment variable
	mongoUri := os.Getenv("MONGODB_URI")
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPass := os.Getenv("MONGODB_PASS")
	// Existence check
	if len(mongoUri) < 1 || len(mongoUser) < 1 || len(mongoPass) < 1 {
		err = errors.New("environment for MONGODB_USER, MONGODB_PASS, MONGODB_URI should be there")
		return mongoUri, err
	}

	strings.Index(mongoUri, "{USER}")
	strings.Index(mongoUri, "{PASS}")

	if !CheckURIFormat(mongoUri) {
		err = errors.New("BAD FORMAT: format of MONGODB_URI should has {USER}:{PASS}. example: mongodb+srv://{USER}:{PASS}@your.database.is.here")
		return mongoUri, err
	}

	mongoUri = strings.Replace(mongoUri, "{USER}", mongoUser, 1)
	mongoUri = strings.Replace(mongoUri, "{PASS}", mongoPass, 1)

	return mongoUri, err
}

func CheckURIFormat(uri_ string) bool {
	return strings.Index(uri_, "{USER}") != -1 && strings.Index(uri_, "{PASS}") != -1
}
