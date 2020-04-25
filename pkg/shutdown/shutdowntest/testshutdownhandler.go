package main

import(
	"com/coder/logging"
	"com/coder/shutdown"
	"github.com/chasex/glog"
)

var Logger *glog.Logger

type Foo struct{
	
}

func (Foo) HandleShutdown(){
	Logger.Infof("HandleShutdown of Foo called....")
}

func init(){
	Logger = logging.GetLogger("test_shutdown", "/tmp", true)
	//Don't set logger if you want to log to the terminal
	//shutdown.GetShutdownHandler().Init(Logger)
	shutdown.GetShutdownHandler().Init(nil)
}

func main() {
	shutdownTest()
	shutdown.Wait()
}

func shutdownTest(){
	foo := &Foo{}
	shutdown.AddListener(foo)
	
}

