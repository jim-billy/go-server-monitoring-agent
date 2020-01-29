package main
import (
	
	"com/coder/scheduler"
	"com/coder/monagent/collector"
	"time"
)

var Logger *glog.Logger

func init(){
	fmt.Println("Init of schedulertest")
	runtime.GOMAXPROCS(3*runtime.NumCPU())
	Logger = logging.GetLogger("schedulertest", "/tmp", true)
}

func testScheduler(){
	var sched scheduler.Scheduler
	var schTask scheduler.ScheduleTask
	schTask.SetName("1 sec schedule")
	schTask.SetType(scheduler.REPETITIVE_TASK)
	schTask.SetInterval(1)
	sched.Schedule(schTask)
}

func testDataCollectionScheduler(){
	var sched scheduler.Scheduler
	sched = scheduler.GetScheduler("DataCollectionScheduler") 
	var schTask scheduler.ScheduleTask
	schTask.SetName("1 sec schedule")
	schTask.SetType(scheduler.REPETITIVE_TASK)
	schTask.SetInterval(1)
	serverMonJob := collector.ServerMonitoringJob{JobType : collector.CPU_STATS, Id : "1"}
	schTask.SetJob(serverMonJob)
	sched.Schedule(schTask)
}

func main(){
	testScheduler()
	//testDataCollectionScheduler()
	time.Sleep(1600 * time.Second)
}