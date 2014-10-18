package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	ApiKey   string      `json:"apiKey"`
	DbConfig MysqlConfig `json:"mysqlConfig"`
	Limits   LimitConfig `json:"limits"`
}

type MysqlConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

type LimitConfig struct {
	RequestsPerTenSeconds int `json:"reqPerTenSec"`
	RequestsPerTenMinutes int `json:"reqPerTenMin"`
}

// Opens configuration files and
// implements associated structs
func openAndReadConfig(configFileName string) (config Configuration) {

	// load config file
	configFile, err := os.Open(configFileName)
	checkErr(err, "Unable to open config file")

	// parse config file
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	checkErr(err, "Unable to decode json")

	return config
}
