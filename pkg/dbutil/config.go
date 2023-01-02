package dbutil

type DBConfig struct {
	DatabaseName        string
	DatabaseURI         string
	MigrationCollection string
}

func NewConfig(dbName_ string, migrationCollection_ string) (DBConfig, error) {
	// Vars
	config := DBConfig{}
	var uri string
	var err error = nil

	// Checking
	if uri, err = GetURI(); err != nil {
		return config, err
	}

	// Set config
	config.DatabaseURI = uri
	config.DatabaseName = dbName_
	config.MigrationCollection = migrationCollection_

	return config, err
}
