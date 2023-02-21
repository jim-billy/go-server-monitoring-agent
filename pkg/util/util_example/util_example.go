package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/routinepool"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/scheduler"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/shutdown"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/util"
)

type WebsiteJob struct {
	Website    string
	Id         int
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

func (websiteJob *WebsiteJob) GetID() int {
	return websiteJob.Id
}

type TestScheduler struct{}

func (TestScheduler) testScheduler() {
	var sched *scheduler.Scheduler
	sched = scheduler.GetScheduler("DataCollectionScheduler")

	var schTask scheduler.ScheduleTask
	schTask.SetName("1 sec schedule")
	schTask.SetType(scheduler.REPETITIVE_TASK)
	schTask.SetInterval(1)
	websiteJob := &WebsiteJob{Website: "https://127.0.0.1/index.html", Id: 1}
	schTask.SetJob(websiteJob)
	sched.Schedule(schTask)
}

func (testScheduler TestScheduler) HandleShutdown() {
	log.Println("HandleShutdown of TestScheduler called....")
	//testScheduler.printSchedulerPerformanceStats()
}

func withError() {
	defer util.CatchPanic(nil, "WorkRoutine", "withError")
	var job *WebsiteJob
	fmt.Println(job.Id)
}

func main() {
	var testSch *TestScheduler
	testSch = new(TestScheduler)
	testSch.testScheduler()
	//Register this scheduler for shutdown notification to call HandleShutdown() by the shutdown module
	shutdown.AddListener(testSch)
	withError()
	shutdown.Wait()
}
