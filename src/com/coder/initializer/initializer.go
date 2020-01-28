package initializer

import (
	"os"
	"fmt"
	"com/coder/util"
)

var agentHome string
var agentLogDir string
var agentConfDir string
var agentDataDir string
var agentScriptsDir string


func SetWorkingDirectories() {
  dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else{
		agentHome = dir
		agentLogDir = agentHome +"/"+ "logs"
		agentConfDir = agentHome +"/"+ "conf"
		agentDataDir = agentHome +"/"+ "data"
		agentScriptsDir = agentHome +"/"+ "scripts"
	}
  fmt.Println("AgentHome : ",agentHome)
  fmt.Println("AgentLogDir : ",agentLogDir)
  fmt.Println("AgentConfDir : ",agentConfDir)
  fmt.Println("AgentDataDir : ",agentDataDir)
  fmt.Println("AgentScriptsDir : ",agentScriptsDir)
  util.CheckAndCreateDirectory(agentLogDir)
  util.CheckAndCreateDirectory(agentDataDir)
}

func GetAgentHome() string{
	return agentHome
}

func GetAgentLogDir() string{
	return agentLogDir
}

func GetAgentConfDir() string{
	return agentConfDir
}

func GetAgentScriptsDir() string{
	return agentScriptsDir
}

func GetAgentDataDir() string{
	return agentDataDir
}

