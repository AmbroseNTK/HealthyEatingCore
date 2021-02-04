package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Address         string `json:"address"`
	ConnectionURL   string `json:"connectionURL"`
	DBName          string `json:"db_name"`
	FirebaseKeyFile string `json:"firebase_key"`
}

func (c *Configuration) Load(configFile string) {
	file, fileErr := os.Open(configFile)
	if fileErr != nil {
		fmt.Println(fileErr)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	if err != nil {
		fmt.Println(err)
	}
}
