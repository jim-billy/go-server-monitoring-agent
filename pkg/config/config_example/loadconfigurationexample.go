package main

import (
	"log"
	"com/coder/config"
)

func main() {
	loadConfiguration()
	
}

func loadConfiguration(){
	filePath := "/tmp/sampleconfig.json"
	appConfig := config.AppConfig{}
	config.LoadConfig(filePath, &appConfig)
	log.Println(appConfig.Db.Host,appConfig.Db.Password)
}