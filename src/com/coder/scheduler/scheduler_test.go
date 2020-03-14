package scheduler 

import (
	"fmt"
	"runtime"
	"time"
	"testing"
	"com/coder/routinepool"
	"com/coder/scheduler"
	"com/coder/logging"
	"com/coder/shutdown"
	"github.com/chasex/glog"
)

var Logger *glog.Logger

type TestScheduler struct{}

func init(){
	fmt.Println("Init of schedulertest")
	runtime.GOMAXPROCS(3*runtime.NumCPU())
	Logger = logging.GetLogger("schedulertest", "/tmp", true)
	//Don't set logger if you want to log to the terminal
	//shutdown.GetShutdownHandler().Init(Logger)
	shutdown.GetShutdownHandler().Init(nil)
}

type WebsiteJob struct {
    Website string
    Id 		int
    ResultData *routinepool.JobResult
}

func (websiteJob *WebsiteJob) DoJob(routinePool *routinepool.RoutinePool) {
    fmt.Println("============================== DoJob : Collecting data for Website : ", websiteJob.Website)
    //time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
    time.Sleep(1 * time.Second)
    /*
    var a[3] int
	j := 5
	fmt.Println(a[j])
	*/
}

func (websiteJob *WebsiteJob) GetId() int{
    return websiteJob.Id
}

func (TestScheduler) testScheduler(){
	var sched *scheduler.Scheduler
	sched = scheduler.GetScheduler("DataCollectionScheduler") 
	
	var schTask scheduler.ScheduleTask
	schTask.SetName("1 sec schedule")
	schTask.SetType(scheduler.REPETITIVE_TASK)
	schTask.SetInterval(1)
	websiteJob := &WebsiteJob{Website : "https://127.0.0.1/index.html", Id : 1}
	schTask.SetJob(websiteJob)
	sched.Schedule(schTask)
}

func (testScheduler TestScheduler) HandleShutdown(){
	Logger.Infof("HandleShutdown of TestScheduler called....")
	//testScheduler.printSchedulerPerformanceStats()
}

func (TestScheduler) printSchedulerPerformanceStats(){
	sched := scheduler.GetScheduler("DataCollectionScheduler")
	sched.PerformanceStats(sched.GetName())
}

func main(){
	
}

func TestGetScheduler( t *testing.T ) {

	var sched *scheduler.Scheduler
	sched = scheduler.GetScheduler("DataCollectionScheduler1") 

	if sched == nil {
		t.Errorf( "scheduler.GetScheduler FAILED, Unable to create scheduler")
	} else {
		t.Logf( "scheduler.GetScheduler PASSED, Scheduler creation successful" )
	}
}

func TestGetScheduler1( t *testing.T ) {

	var testSch *TestScheduler
	testSch = new(TestScheduler)
	testSch.testScheduler()
	//Register this scheduler for shutdown notification to call HandleShutdown() by the shutdown module
	shutdown.AddListener(testSch)
	shutdown.Wait()
}
