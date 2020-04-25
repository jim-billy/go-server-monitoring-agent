package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/chasex/glog"
)

var loggerMap map[string]*glog.Logger

func init() {
	loggerMap = make(map[string]*glog.Logger)
}

func deleteLog(deleteFileName string, logDir string) {
	//fmt.Println("Inside Delete log")
	files, err := filepath.Glob(logDir + "/" + deleteFileName + "*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println("Deleting file : ", f)
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func createLogger(fname string, logDir string) *glog.Logger {
	fmt.Println("Creating Logger : ", fname)
	if logDir == "" {
		logDir = "/tmp"
	}
	options := glog.LogOptions{
		File:    logDir + "/" + fname + ".log",
		Flag:    glog.LstdFlags,
		Level:   glog.Ldebug,
		Mode:    glog.R_Size,
		Maxsize: 5 * 1024 * 1024,
	}
	var err error
	var logger *glog.Logger
	logger, err = glog.New(options)
	if err != nil {
		panic(err)
		//fmt.Println("createLogger error..")
	} else {
		loggerMap[fname] = logger
	}
	//fmt.Printf("End of createLogger : %v \n",logger)
	return logger
}

func logBasicInfo(logger *glog.Logger) {
	logger.Infof("========================== Logging started ==========================")
	logger.Infof("No. of CPU cores : " + strconv.Itoa(runtime.NumCPU()))

}

//Public methods

func GetLogger(loggerName string, logDir string, deleteOldLogs bool) *glog.Logger {
	loggerToReturn := loggerMap[loggerName]
	if loggerToReturn == nil && logDir != "" {
		if deleteOldLogs {
			deleteLog(loggerName, logDir)
		}
		loggerToReturn = createLogger(loggerName, logDir)
		logBasicInfo(loggerToReturn)
		loggerToReturn.Flush()
	}
	return loggerToReturn
}

func FlushAllLogs() {
	for logName, logger := range loggerMap {
		logger.Infof("Flusing logger : " + logName)
		logger.Flush()

	}
}
