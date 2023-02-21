package main

import (
	"fmt"
	"log"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/logging"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/shutdown"
)

var Logger *log.Logger

type Foo struct {
}

func (Foo) HandleShutdown() {
	fmt.Println("HandleShutdown of Foo called....")
}

func init() {
	Logger = logging.GetLogger("test_shutdown", "/tmp", true)
	//Don't set logger if you want to log to the terminal
	//shutdown.GetShutdownHandler().Init(Logger)
	shutdown.GetShutdownHandler().Init(nil)
}

func main() {
	shutdownTest()
	shutdown.Wait()
}

func shutdownTest() {
	foo := &Foo{}
	shutdown.AddListener(foo)

}
