package main

import (
	"com/coder/monagent/agentconstants"
	"com/coder/logging"
	"com/coder/initializer"
	"com/coder/util"
	"com/coder/shutdown"
	"com/coder/monagent/collector"
	"fmt"
	"runtime"
)

func init(){
	fmt.Println("Init of MonAgent")
	initializer.SetWorkingDirectories()
	initLogging()
	agentconstants.Initialize()
	loadConfiguration()
	runtime.GOMAXPROCS(3*runtime.NumCPU())
	shutdown.GetShutdownHandler().Init(nil)
}

func initLogging(){
	fmt.Println("=================== initLogging ============== ")
	agentconstants.Logger = logging.GetLogger("agent", initializer.GetAgentLogDir(), true)
	agentconstants.ErrorLogger = logging.GetLogger("error", initializer.GetAgentLogDir(), true)
}

func loadConfiguration(){
	collector.GetcollectorAPI().Initialize()
}

func main() {
	util.SetLimit()
	collector.GetcollectorAPI().ScheduleDataCollection()
	shutdown.Wait()
}


