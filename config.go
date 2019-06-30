package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var config = "./config.json"
var defaultConfig = `{
  "port": ":8090"
}`

type Config struct {
	Port string `json:"port"`
}

func readConfig() (Config, error) {
	data, err := ioutil.ReadFile(config)
	if err != nil {
		return Config{}, err
	}
	conf := Config{}
	err = json.Unmarshal(data, &conf)
	return conf, err
}

func writeConfig() {
	err := ioutil.WriteFile(config, []byte(defaultConfig), 0644)
	if err != nil {
		log.Printf("config.json created error, %s", err.Error())
	} else {
		log.Printf("config.json had be created, please check and restart server.")
	}
}
