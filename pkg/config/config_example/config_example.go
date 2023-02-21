package main

import (
	"encoding/json"
	"log"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/config"
)

func main() {
	loadConfiguration()

}

type AppConfig struct {
	Name         string               `json:"name"`
	Script       bool                 `json:"script"`
	Command      string               `json:"command"`
	CommandArgs  string               `json:"command_args"`
	Interval     int                  `json:"interval"`
	KeyValue     bool                 `json:"key_value"`
	ParseAllLine bool                 `json:"parse_all_line"`
	ParseImpl    string               `json:"parse_impl"`
	ParseConfig  []ParseConfiguration `json:"parse_config"`
}

// ParseConfiguration holds the attributes defined in the linux_monitors.json that are necessary for parsing the collected data.
type ParseConfiguration struct {
	MetricName string `json:"metric_name"`
	ParseLine  int    `json:"parse_line"`
	Delimiter  string `json:"delimiter"`
	Token      int    `json:"token"`
	Counter    bool   `json:"counter"`
	Expression string `json:"expression"`
}

func loadConfiguration() {
	var appConfig AppConfig
	filePath := "./sampleconfig.json"
	configLoader := config.GetConfigLoader()
	byteArr, _ := configLoader.LoadBytesFromJson(filePath)
	json.Unmarshal(byteArr, &appConfig)

	log.Println(appConfig.Name, appConfig.Script)
}
