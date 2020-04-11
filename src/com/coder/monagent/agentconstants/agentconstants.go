package agentconstants

import (
	"com/coder/initializer"

	"github.com/chasex/glog"
)

const (
	// AgentName denotes the name of the agent
	AgentName = "ServerMonitoringAgent"
	// AgentConfFileName denotes the configuration file name of the agent
	AgentConfFileName = "agent_conf.json"
	// LinuxMonitorsFileName denotes the linux monitors file name of the agent
	LinuxMonitorsFileName = "linux_monitors.json"
)

// Logger stores the instance of the logger
var Logger *glog.Logger

// ErrorLogger stores the instance of the error logger
var ErrorLogger *glog.Logger

// AgentConfFilePath denotes the configuration file path
var AgentConfFilePath string

// LinuxMonitorsFilePath denotes the linux monitors file path
var LinuxMonitorsFilePath string

// Initialize is used for setting agent constants.
func Initialize() {
	AgentConfFilePath = initializer.GetAgentConfDir() + "/" + AgentConfFileName
	LinuxMonitorsFilePath = initializer.GetAgentConfDir() + "/" + LinuxMonitorsFileName
}
