package agentconstants

import (
	"com/coder/initializer"
)

const AGENT_NAME = "MonitoringAgent"
const AGENT_CONF_FILE_NAME = "agent_conf.json"
const LINUX_MONITORS_FILE_NAME = "linux_monitors.json"

var AGENT_CONF_FILE_PATH string
var LINUX_MONITORS_FILE_PATH string

func Initialize(){
	AGENT_CONF_FILE_PATH = initializer.GetAgentConfDir() + "/" + AGENT_CONF_FILE_NAME
	LINUX_MONITORS_FILE_PATH = initializer.GetAgentConfDir() + "/" + LINUX_MONITORS_FILE_NAME 
}

