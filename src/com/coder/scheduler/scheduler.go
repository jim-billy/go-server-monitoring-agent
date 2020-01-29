package scheduler

import(
	"fmt"
	"time"
	"com/coder/routinepool"
)

const(
	REPETITIVE_TASK="REPETITIVE_TASK"
	ONE_TIME_TASK="ONE_TIME_TASK"
	DEFAULT_NO_OF_WORKERS = 8
	QUEUE_CAPACITY = 100
)


var schedulerMap map[string]Scheduler

func init(){
	schedulerMap = make(map[string]Scheduler)
}

type ScheduleTask struct {  
	name string
	taskType string
	interval int
	scheduleJob routinepool.Job
}

func (schTask *ScheduleTask) SetName(name string){
	schTask.name = name
}

func (schTask *ScheduleTask) GetName() string{
	return schTask.name
}

func (schTask *ScheduleTask) SetType(taskType string){
	schTask.taskType = taskType
}

func (schTask *ScheduleTask) GetType() string{
	return schTask.taskType
}

func (schTask *ScheduleTask) SetInterval(interval int){
	schTask.interval = interval
}

func (schTask *ScheduleTask) GetInterval() int{
	return schTask.interval
}

func (schTask *ScheduleTask) SetJob(schJob routinepool.Job){
	schTask.scheduleJob = schJob
}

func (schTask *ScheduleTask) GetJob() routinepool.Job{
	return schTask.scheduleJob
}

type Scheduler struct {  
	name string
	routinePool *routinepool.RoutinePool
}

func (sch *Scheduler) Schedule(schTask ScheduleTask){
	if(schTask.taskType == REPETITIVE_TASK){
		fmt.Println("REPETITIVE_TASK")
		ticker := time.NewTicker(time.Duration(schTask.GetInterval()) * time.Second)
	    done := make(chan bool)
	
	    go func() {
	        for {
	            select {
	            case <-done:
	                return
	            case t := <-ticker.C:
					fmt.Println("============================ Sending the job to the worker : ", schTask.scheduleJob,t)
					sch.routinePool.GetJobChannel() <- schTask.scheduleJob
	            }
	        }
	    }()
	}else if(schTask.taskType == ONE_TIME_TASK){
		fmt.Println("ONE_TIME_TASK")
		fmt.Println("============================ Sending the job to the worker : ", schTask.scheduleJob)
		sch.routinePool.GetJobChannel() <- schTask.scheduleJob
	}else{
		fmt.Println("Unknown task type in the input ScheduleTask : ",schTask)
	}
	
}

//Public methods

func GetScheduler(schedulerName string) *Scheduler {
	var toReturn *Scheduler
	for name, scheduler := range schedulerMap { 
	    //fmt.Printf("key[%s] value[%s]\n", name, Scheduler)
	    if name == schedulerName{
	    	toReturn = scheduler
	    }
	}
	if(scheduler == nil){
		fmt.Println("Creating new scheduler : ", schedulerName)
		config := routinepool.RoutinePoolConfig{
			RoutinePoolName: "SchedulerDataCollectionPool",
			RoutinePoolSize: DEFAULT_NO_OF_WORKERS,
			QueueCapacity: QUEUE_CAPACITY,
			//Logger: Logger,
		}

		routinePool, err := routinepool.NewRoutinePool(config)
		if err != nil {
			panic(err)
			fmt.Println("Error while creating worker pool for the scheduler ..",schedulerName)
		}else{
			toReturn.name = schedulerName
			toReturn.routinePool = routinePool
		}
	}
	return toReturn
}


