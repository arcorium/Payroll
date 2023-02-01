package dbutil

import (
	"Penggajian/util"
	"fmt"
	"os"
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
		return config, err
	}

	// Get binding address to listen
	bindAddr := bindIp + ":" + bindPort

	// Get mongodb uri
	if uri, err = GetURI(); err != nil {
		return config, err
	}

	// Set config
	config.DatabaseURI = uri
	config.DatabaseName = dbName_
	config.MigrationCollection = migrationCollection_
	config.BindAddress = bindAddr
	config.SecretKey = secretKey

	fmt.Println(config.DatabaseURI)

	return config, err
}
