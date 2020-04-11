package main

import (
	"com/coder/initializer"
	"com/coder/logging"
	"com/coder/monagent/agentconstants"
	"com/coder/monagent/collector"
	"com/coder/shutdown"
	"com/coder/util"
	"fmt"
	"runtime"
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
