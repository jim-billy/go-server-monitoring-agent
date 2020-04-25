package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"fmt"
	"io/ioutil"
)

var configLoader ConfigLoader

type ConfigLoader struct{}

func LoadConfig(filePath string) *json.Decoder{
	c := flag.String("c", filePath, "Specify the configuration file.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder
}

func (configLoader *ConfigLoader) LoadBytesFromJson(filePath string) ([]byte, error){
	var byteArr []byte
	var errToReturn error
	c := flag.String("c", filePath, "Specify the configuration file.")
	flag.Parse()
    jsonFile, errToReturn := os.Open(*c)
    
    if errToReturn != nil {
        fmt.Println("can't open config file: ",errToReturn)
        return byteArr, errToReturn
    }else{
	    fmt.Println("Successfully read configuration file : ",filePath)	
    }
    
    defer jsonFile.Close()
    
    byteArr, errToReturn = ioutil.ReadAll(jsonFile)
    return byteArr, errToReturn
    
}

func GetConfigLoader() *ConfigLoader{
	return &configLoader
}

