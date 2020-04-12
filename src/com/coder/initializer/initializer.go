package initializer

import (
	"com/coder/util"
	"fmt"
	"os"
)

var agentHome string
var agentLogDir string
var agentConfDir string
var agentDataDir string
var agentScriptsDir string

// SetWorkingDirectories initializes the working directories of the agent
func SetWorkingDirectories() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		agentHome = dir
		agentLogDir = agentHome + "/" + "logs"
		agentConfDir = agentHome + "/" + "conf"
		agentDataDir = agentHome + "/" + "data"
		agentScriptsDir = agentHome + "/" + "scripts"
	}
	fmt.Println("AgentHome : ", agentHome)
	fmt.Println("AgentLogDir : ", agentLogDir)
	fmt.Println("AgentConfDir : ", agentConfDir)
	fmt.Println("AgentDataDir : ", agentDataDir)
	fmt.Println("AgentScriptsDir : ", agentScriptsDir)
	util.CheckAndCreateDirectory(agentLogDir)
	util.CheckAndCreateDirectory(agentDataDir)
}

// GetAgentHome returns agent's home directory
func GetAgentHome() string {
	return agentHome
}

// GetAgentLogDir returns agent's log directory
func GetAgentLogDir() string {
	return agentLogDir
}

// GetAgentConfDir returns agent's conf directory
func GetAgentConfDir() string {
	return agentConfDir
}

// GetAgentScriptsDir returns agent's scripts directory
func GetAgentScriptsDir() string {
	return agentScriptsDir
}

// GetAgentDataDir returns agent's data directory
func GetAgentDataDir() string {
	return agentDataDir
}
