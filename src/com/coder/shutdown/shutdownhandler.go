package shutdown

import (
	"fmt"
	"com/coder/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/chasex/glog"
)

var shutdownHanler ShutdownHandler

type ShutdownHandler struct{
	signalChan chan os.Signal
	cleanupDone chan bool
	listener []ShutdownListener	
	logger *glog.Logger
}

type ShutdownListener interface{
	HandleShutdown()
}
/*
	Don't set logger if you want to log to the terminal
	shutdown.GetShutdownHandler().Init(nil)
*/
func (shutdownHandler *ShutdownHandler) Init(logger *glog.Logger){
	shutdownHanler.logger = logger
	shutdownHanler.log("Init of ShutdownHandler")
	shutdownHanler.signalChan = make(chan os.Signal, 1)
	shutdownHanler.cleanupDone = make(chan bool, 1)
	signal.Notify(shutdownHanler.signalChan, syscall.SIGINT, syscall.SIGTERM)
	shutdownHanler.registerForShutdown()
}

func (shutdownHandler *ShutdownHandler) registerForShutdown(){
	go func() {
		sig := <-shutdownHanler.signalChan
        shutdownHanler.log("Shutdown signal received : "+sig.String()+". Stopping all modules")
	    shutdownHandler.notifyListeners()
	    logging.FlushAllLogs()
	    time.Sleep(time.Millisecond * 500)
	    shutdownHanler.cleanupDone <- true
	}()
}

func (shutdownHandler *ShutdownHandler) notifyListeners(){
	for i := 0; i < len(shutdownHandler.listener); i++ {
        shutdownHandler.listener[i].HandleShutdown()
    }
}

func (shutdownHandler *ShutdownHandler) log(message string){
	if(shutdownHandler.logger == nil){
		fmt.Println(message)	
	}else{
		shutdownHandler.logger.Infof(message)	
	}
}

func AddListener(listener ShutdownListener){
	shutdownHanler.listener = append(shutdownHanler.listener, listener)
}

func GetShutdownHandler() *ShutdownHandler{
	return &shutdownHanler
}

/*
Blocks the calling thread and waits till the shutdown notification is received
*/
func Wait(){
	shutdownHanler.log("Listening for shutdown signal")
    <-shutdownHanler.cleanupDone
    os.Exit(3)
}