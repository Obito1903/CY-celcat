package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	UserName     string   `json:"userName"`
	UserPassword string   `json:"userPassword"`
	CelcatHost   string   `json:"celcatHost"`
	Groupes      []Groupe `json:"groupes"`
}

type Groupe struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func ReadConfig(path string) Config {
	configJson, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal("Couldn't open config file", err)
	}
	// Wait for the file to be fully read
	defer configJson.Close()
	// Convert file int a byte field
	configByte, _ := ioutil.ReadAll(configJson)

	var config Config
	json.Unmarshal(configByte, &config)
	return config
}
