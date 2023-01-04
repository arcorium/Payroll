package dbutil

import (
	"Penggajian/pkg/util"
	"fmt"
	"os"
)

type DBConfig struct {
	DatabaseName        string
	DatabaseURI         string
	MigrationCollection string
	BindAddress         string
}

func NewConfig(dbName_ string, migrationCollection_ string) (DBConfig, error) {
	// Vars
	config := DBConfig{}
	var uri string
	var err error = nil

	bindIp := os.Getenv("BIND_IP")
	bindPort := os.Getenv("BIND_PORT")

	if util.IsEmpty(bindIp) || util.IsEmpty(bindPort) {
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

	fmt.Println(config.DatabaseURI)

	return config, err
}
