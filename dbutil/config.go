package dbutil

import (
	"errors"
	"log"
	"os"

	"Penggajian/util"
)

type DBConfig struct {
	DatabaseName        string
	DatabaseURI         string
	MigrationCollection string
	BindAddress         string
	SecretKey           string
}

func NewConfig(dbName_ string, migrationCollection_ string) (DBConfig, error) {
	// Vars
	config := DBConfig{}
	var uri string
	var err error = nil

	bindIp := os.Getenv("BIND_IP")
	bindPort := os.Getenv("BIND_PORT")
	secretKey := os.Getenv("SECRET_KEY")

	if util.IsEmpty(bindPort) || util.IsEmpty(secretKey) {
		bindPort = os.Getenv("PORT")
		if util.IsEmpty(bindPort) {
			bindPort = "3000"
		}
		return config, errors.New("environment is not satisfied")
	}

	// Get binding address to listen
	bindAddr := bindIp + ":" + bindPort

	// Get mongodb uri
	//if uri, err = GetURI(); err != nil {
	//	return config, err
	//}

	// Set config
	//config.DatabaseURI = uri
	uri = os.Getenv("MONGODB_URI")
	if util.IsEmpty(uri) {
		return config, errors.New("environment is not satisfied")
	}
	config.DatabaseURI = uri
	config.DatabaseName = dbName_
	config.MigrationCollection = migrationCollection_
	config.BindAddress = bindAddr
	config.SecretKey = secretKey

	log.Println("Config Database:", config.DatabaseURI)
	log.Println("Config Database Name:", config.DatabaseName)
	log.Println("Config Listen:", config.BindAddress)

	return config, err
}
