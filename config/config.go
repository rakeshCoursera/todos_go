package config

import (
    "encoding/json"
    "fmt"
    "os"
)

// configuration file
type Config struct {
	Port                int     `json:"port"`
	Region							string	`json:"region"`
	Access_Key          string  `json:"access_key"`
	Secret_Key          string  `json:"secret_key"`
	Access_Token        string  `json:"access_token"`
	Table_Name          string  `json:"table_name"`
}

func LoadConfiguration() Config {
	var config Config
	configFile, err := os.Open("config/config.json")
	defer configFile.Close()
	if err != nil {
			fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
