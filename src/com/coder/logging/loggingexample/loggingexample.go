package main

import(
	"com/coder/logging"
	"github.com/chasex/glog"
	"fmt"
	"strconv"
	"runtime"
)

var Logger1 *glog.Logger
var Logger2 *glog.Logger
var Logger3 *glog.Logger

func init(){
	fmt.Println("Init of logging_example")
	Logger1 = logging.GetLogger("logging1")
	Logger2 = logging.GetLogger("logging2")
	Logger3 = logging.GetLogger("logging3")
	printMessages("logging1")
	printMessages("logging2")
	printMessages("logging3")
}

func printMessages(logName string){
	logger := logging.GetLogger(logName)
	logger.Infof("========================== Logging started =========================="+logName)
	logger.Infof("No. of CPU cores : "+strconv.Itoa(runtime.NumCPU()))	
}

func main(){
	logging.FlushAllLogs()
}