// Package for scheduling jobs at specified time.
package scheduler

import (
	//"fmt"
	"time"

	"github.com/chasex/glog"
	"github.com/gojavacoder/go-server-monitoring-agent/pkg/routinepool"
)

type ScheduleTaskType int

const (
	ONE_TIME_TASK         ScheduleTaskType = 1
	REPETITIVE_TASK       ScheduleTaskType = 2
	DEFAULT_NO_OF_WORKERS                  = 8
	QUEUE_CAPACITY                         = 100
)

var schedulerMap map[string]*Scheduler

func init() {
	schedulerMap = make(map[string]*Scheduler)
}

func (taskType ScheduleTaskType) Get() ScheduleTaskType {
	return taskType
}

type ScheduleTask struct {
	name        string
	taskType    ScheduleTaskType
	interval    int
	scheduleJob routinepool.Job
}

func (schTask *ScheduleTask) SetName(name string) {
	schTask.name = name
}

func (schTask *ScheduleTask) GetName() string {
	return schTask.name
}

func (schTask *ScheduleTask) SetType(taskType ScheduleTaskType) {
	schTask.taskType = taskType
}

func (schTask *ScheduleTask) GetType() ScheduleTaskType {
	return schTask.taskType
}

func (schTask *ScheduleTask) SetInterval(interval int) {
	schTask.interval = interval
}

func (schTask *ScheduleTask) GetInterval() int {
	return schTask.interval
}

func (schTask *ScheduleTask) SetJob(schJob routinepool.Job) {
	schTask.scheduleJob = schJob
}

func (schTask *ScheduleTask) GetJob() routinepool.Job {
	return schTask.scheduleJob
}

type Scheduler struct {
	name        string
	routinePool *routinepool.RoutinePool
	logger      *glog.Logger
}

func (sch *Scheduler) SetName(name string) {
	sch.name = name
}

func (sch *Scheduler) GetName() string {
	return sch.name
}

func (sch *Scheduler) SetLogger(logger *glog.Logger) {
	sch.logger = logger
	sch.routinePool.SetLogger(logger)
}

func (sch *Scheduler) GetLogger() *glog.Logger {
	return sch.logger
}

func (sch *Scheduler) Schedule(schTask ScheduleTask) {
	if schTask.taskType == REPETITIVE_TASK {
		ticker := time.NewTicker(time.Duration(schTask.GetInterval()) * time.Second)
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					sch.logger.Infof("Scheduler : Schedule : ============================ Sending the job to the worker : %v %v", schTask.scheduleJob, t)
					sch.routinePool.ExecuteJob(schTask.scheduleJob)
				}
			}
		}()
	} else if schTask.taskType == ONE_TIME_TASK {
		sch.logger.Infof("Scheduler : Schedule : ONE_TIME_TASK : ============================ Sending the job to the worker : %v", schTask.scheduleJob)
		sch.routinePool.ExecuteJob(schTask.scheduleJob)
	} else {
		sch.logger.Infof("Scheduler : Schedule : Unknown task type in the input ScheduleTask : %v", schTask)
	}

}

func (sch *Scheduler) PerformanceStats(schedulerName string) {
	routinePool := routinepool.GetRoutinePool(schedulerName + "-RoutinePool")
	routinePool.PerformanceStats()
}

//Public methods

func GetScheduler(schedulerName string) *Scheduler {
	var schToReturn *Scheduler
	for name, scheduler := range schedulerMap {
		//fmt.Printf("key[%s] value[%s]\n", name, Scheduler)
		if name == schedulerName {
			schToReturn = scheduler
		}
	}
	if schToReturn == nil {
		config := routinepool.RoutinePoolConfig{
			RoutinePoolName: schedulerName + "-RoutinePool",
			RoutinePoolSize: DEFAULT_NO_OF_WORKERS,
			QueueCapacity:   QUEUE_CAPACITY,
		}

		routinePool, err := routinepool.NewRoutinePool(config)
		if err != nil {
			panic(err)
		} else {
			schToReturn = new(Scheduler)
			schToReturn.name = schedulerName
			schToReturn.routinePool = routinePool
		}
	}
	return schToReturn
}
