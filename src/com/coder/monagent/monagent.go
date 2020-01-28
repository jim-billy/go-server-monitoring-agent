package main
import (
	//"time"
	"com/coder/monagent/agentconstants"
	"com/coder/logging"
	"com/coder/initializer"
	"com/coder/util"
	"com/coder/shutdown"
	"com/coder/monagent/collector"
	"fmt"
	"runtime"
	"github.com/chasex/glog"
)


var Logger *glog.Logger
var ErrorUrlLogger *glog.Logger
var SuccessUrlLogger *glog.Logger


func init(){
	fmt.Println("Init of MonAgent")
	initializer.SetWorkingDirectories()
	initLogging()
	agentconstants.Initialize()
	loadConfiguration()
	runtime.GOMAXPROCS(3*runtime.NumCPU())
}

func initLogging(){
	fmt.Println("=================== initLogging ============== ")
	Logger = logging.GetLogger("agent", initializer.GetAgentLogDir())
	ErrorUrlLogger = logging.GetLogger("error", initializer.GetAgentLogDir())
	SuccessUrlLogger = logging.GetLogger("success", initializer.GetAgentLogDir())
}

func loadConfiguration(){
	collector.GetCollectorApi().Initialize()
}

func main() {
	util.SetLimit()
	collector.GetCollectorApi().ScheduleDataCollection()
	shutdown.ShutdownListener()
}


