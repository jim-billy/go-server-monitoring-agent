package main

import (
	"fmt"
	"runtime"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/initializer"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/logging"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/monagent/agentconstants"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/monagent/collector"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/shutdown"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/util"
)

// Initializes the server monitoring agent
func init() {
	fmt.Println("Init of MonAgent")
	initializer.SetWorkingDirectories()
	initLogging()
	agentconstants.Initialize()
	loadConfiguration()
	runtime.GOMAXPROCS(3 * runtime.NumCPU())
	shutdown.GetShutdownHandler().Init(nil)
}

//Initializes logging
func initLogging() {
	fmt.Println("=================== initLogging ============== ")
	agentconstants.Logger = logging.GetLogger("agent", initializer.GetAgentLogDir(), true)
	agentconstants.ErrorLogger = logging.GetLogger("error", initializer.GetAgentLogDir(), true)
}

// Loads the configuration needed for data collection by initializing the CollectorAPI
func loadConfiguration() {
	collector.GetcollectorAPI().Initialize()
}

// Entry point for the server monitoring agent
func main() {
	util.SetLimit()
	collector.GetcollectorAPI().ScheduleDataCollection()
	shutdown.Wait()
}
