package routinepool

import(
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
	"errors"
	"com/coder/util"
	"com/coder/shutdown"
	"github.com/chasex/glog"
)

const (
	MAX_ROUTINE_POOL_SIZE int = 500
	MAX_QUEUE_CAPACITY int = 1000
)

var routinePoolMap map[string]*RoutinePool

func init(){
	routinePoolMap = make(map[string]*RoutinePool)
}

type Job interface {
	DoJob(routinePool *RoutinePool)
	GetId() int
}
	
type JobResult struct {
	Result map[string]interface{}
}

type RoutinePoolConfig struct {
	RoutinePoolName 	string
	RoutinePoolSize    	int
	QueueCapacity    	int //Maximum number of jobs that can be added to the routinepool without blocking the calling thread.
	Logger *glog.Logger
}

func (poolConfig RoutinePoolConfig) String() string {
    return fmt.Sprintf("RoutinePoolName : %s, RoutinePoolSize : %d, QueueCapacity : %d", poolConfig.RoutinePoolName, poolConfig.RoutinePoolSize, poolConfig.QueueCapacity)
}

type RoutinePool struct {
	poolConfig    RoutinePoolConfig
	shutdownChannel  	 chan bool //Channel used to shut down the work routines.
	jobChannel           chan Job //Channel used to process the incoming jobs
	completedJobsChannel chan Job
	resultChannel chan Job
	queuedWork           int32 //The number of work items queued.
	activeRoutines       int32 //The number of routines active.
	workRoutines map[string]*WorkRoutine
}

func NewRoutinePool(routinePoolConfig RoutinePoolConfig) (*RoutinePool, error){
	//Evaluate the input config before constructing the routinepool
	err := evaluateRoutinePoolConfig(routinePoolConfig)
	if(err != nil){
		return nil, err;
	}
	routinePool := &RoutinePool{
		poolConfig: routinePoolConfig,
		shutdownChannel:      make(chan bool),
		//All the below channels are non-blocking until the routinePoolConfig.QueueCapacity is reached
		jobChannel:           make(chan Job, routinePoolConfig.QueueCapacity),
		completedJobsChannel: make(chan Job, routinePoolConfig.QueueCapacity),
		resultChannel: 		  make(chan Job, routinePoolConfig.QueueCapacity),
		queuedWork:           0,
		activeRoutines:       0,
		workRoutines:		  make(map[string]*WorkRoutine),
	}
	routinePoolMap[routinePoolConfig.RoutinePoolName] = routinePool
	shutdown.AddListener(routinePool)
	for i := 1; i <= routinePoolConfig.RoutinePoolSize; i++ {
		workRoutine := NewWorkRoutine(routinePool, i)
		go workRoutine.run()
	}
	return routinePool, nil
}

//Private method for evaluating input RoutinePoolConfig
func evaluateRoutinePoolConfig(routinePoolConfig RoutinePoolConfig) error{
	if(routinePoolConfig.RoutinePoolName == ""){
		return errors.New("RoutinePoolName cannot be empty")
	}
	if _, isAlreadyPresent := routinePoolMap[routinePoolConfig.RoutinePoolName]; isAlreadyPresent {
	    return errors.New("RoutinePool with the name '"+routinePoolConfig.RoutinePoolName+"' already exists. Please provide a different name to uniquely identify the RoutinePool")
	}
	if(routinePoolConfig.RoutinePoolSize < 0 || routinePoolConfig.RoutinePoolSize > MAX_ROUTINE_POOL_SIZE){
		return errors.New("RoutinePoolSize should be greater than zero and less than the MAX_ROUTINE_POOL_SIZE : "+strconv.Itoa(MAX_ROUTINE_POOL_SIZE))
	}
	if(routinePoolConfig.QueueCapacity < 0 || routinePoolConfig.QueueCapacity > MAX_QUEUE_CAPACITY){
		return errors.New("QueueCapacity should be greater than zero and less than the MAX_QUEUE_CAPACITY : "+strconv.Itoa(MAX_QUEUE_CAPACITY))	
	}
	return nil
}

func (routPool *RoutinePool) GetLogger() *glog.Logger{
	return routPool.poolConfig.Logger
}

func (routPool *RoutinePool) log(message string){
	logger := routPool.GetLogger()
	//strMessage := fmt.Sprintf("",message...)
	if(logger == nil){
		fmt.Println(message)	
	}else{
		logger.Infof(message)	
	}
}

func (routPool *RoutinePool) ExecuteJob(job Job) bool{
	routPool.incrementQueuedWork()
	routPool.jobChannel <- job
	return true
}

func (routPool *RoutinePool) GetShutdownChannel() chan bool {
	return routPool.shutdownChannel
}

func (routPool *RoutinePool) GetCompletedJobsChannel() chan Job {
	return routPool.completedJobsChannel
}

func (routPool *RoutinePool) GetResultChannel() chan Job {
	return routPool.resultChannel
}

// GetQueuedWork will return the number of work items in queue.
func (routPool *RoutinePool) GetQueuedWork() int32 {
	return atomic.AddInt32(&routPool.queuedWork, 0)
}

func (routinePool *RoutinePool) incrementQueuedWork(){
	atomic.AddInt32(&routinePool.queuedWork, 1)	
}

func (routinePool *RoutinePool) decrementQueuedWork(){
	atomic.AddInt32(&routinePool.queuedWork, -1)	
}

// GetActiveRoutines will return the number of routines performing work.
func (routPool *RoutinePool) GetActiveRoutines() int32 {
	//fmt.Printf("POINTER :: The address of the received routPool in ActiveRoutines : %p\n", routPool)
	return atomic.AddInt32(&routPool.activeRoutines, 0)
}

func (routinePool *RoutinePool) incrementActiveRoutines(){
	atomic.AddInt32(&routinePool.activeRoutines, 1)	
}

func (routinePool *RoutinePool) decrementActiveRoutines(){
	atomic.AddInt32(&routinePool.activeRoutines, -1)	
}

type WorkRoutine struct {
	workRoutineName string
	startTime time.Time
	totalJobs int
	maxJobTime int64
	minJobTime int64
	routinePool *RoutinePool
	
}

func NewWorkRoutine(routinePool *RoutinePool, routineId int) *WorkRoutine{ 
	workRoutine := &WorkRoutine{
		workRoutineName: routinePool.poolConfig.RoutinePoolName+"-WorkRoutine-"+strconv.Itoa(routineId),
		routinePool: routinePool,
	}
	routinePool.workRoutines[workRoutine.workRoutineName] = workRoutine
	return workRoutine
}

func (workRoutine WorkRoutine) String() string {
    return fmt.Sprintf("WorkRoutine : %s, Total jobs : %s, Max job time : %s, Min job time : %s", workRoutine.workRoutineName,strconv.Itoa(workRoutine.totalJobs),strconv.FormatInt(workRoutine.maxJobTime, 10),strconv.FormatInt(workRoutine.minJobTime,10))
}

func (workRoutine *WorkRoutine) run() {
	WORK_ROUTINE_LABEL :
	    for {
	        select {
	        	case shutdown := <-workRoutine.routinePool.shutdownChannel:
		        	workRoutine.routinePool.log("Shutting down work routine : "+workRoutine.workRoutineName+" Shutdown received : "+ strconv.FormatBool(shutdown))
		            //fmt.Println(workRoutine.routinePool.poolConfig.RoutinePoolName+", worker : "+strconv.Itoa(workRoutine)+" : Shutdown received : ", shutdown)
		            break WORK_ROUTINE_LABEL
		        case job := <- workRoutine.routinePool.jobChannel:
			        //fmt.Printf("POINTER :: The address of the received routinePool in workRoutine : %p\n", routinePool)
		        	//routinePool.poolConfig.RoutinePoolLogger.Infof(routinePool.poolConfig.RoutinePoolName+", worker : "+strconv.Itoa(workRoutine)+" : received message : ", job)
		        	//routinePool.poolConfig.RoutinePoolLogger.Infof("Routine "+routinePool.poolConfig.RoutinePoolName+" received message : "+msg1)
		        	//fmt.Println("??????????????????????????????? workRoutine Id : ",strconv.Itoa(workRoutine),", JobId : ",strconv.Itoa(job.GetId()))
		        	workRoutine.startJobTime()
		        	workRoutine.incrementTotalJobs()
		        	workRoutine.routinePool.log("workRoutine Id : : "+workRoutine.workRoutineName+", JobId : "+strconv.Itoa(job.GetId()))
		        	workRoutine.safelyDoWork(job)
		        	workRoutine.endJobTime()
	        }
	    }
}

func (workRoutine *WorkRoutine) safelyDoWork(job Job) {
	defer util.CatchPanic(nil, "WorkRoutine", "SafelyDoWork")
	defer workRoutine.routinePool.decrementActiveRoutines()
	
	workRoutine.routinePool.decrementQueuedWork()
	workRoutine.routinePool.incrementActiveRoutines()
	
	job.DoJob(workRoutine.routinePool)
}

func (workRoutine *WorkRoutine) incrementTotalJobs() {
	workRoutine.totalJobs++
}

func (workRoutine *WorkRoutine) startJobTime() {
	workRoutine.startTime = time.Now()
}

func (workRoutine *WorkRoutine) endJobTime() {
	//elapsed := time.Since(workRoutine.startTime).Milliseconds()
	var elapsed int64 = time.Since(workRoutine.startTime).Nanoseconds() / 1e6
	if(workRoutine.maxJobTime == 0){
		workRoutine.maxJobTime = elapsed
		workRoutine.minJobTime = elapsed
		return
	}
	if(workRoutine.maxJobTime < elapsed){
		workRoutine.maxJobTime = elapsed
	}
	if(workRoutine.minJobTime > elapsed){
		workRoutine.minJobTime = elapsed
	}
}

func (routinePool *RoutinePool) GetStats() {
	routinePool.log("QueuedWork : "+strconv.Itoa(int(routinePool.GetQueuedWork())))
	routinePool.log("ActiveRoutines : "+strconv.Itoa(int(routinePool.GetActiveRoutines())))
	
	for _, workRoutine := range routinePool.workRoutines { 
	    routinePool.log(workRoutine.String())
	}    
}

func (*RoutinePool) HandleShutdown() {
	ShutdownRoutinePools()
}

//Public methods

func GetRoutinePool(routinePoolName string) *RoutinePool {
	var toReturn *RoutinePool
	for name, routinePool := range routinePoolMap { 
	    //fmt.Printf("key[%s] value[%s]\n", name, routinePool)
	    if name == routinePoolName{
	    	toReturn = routinePool
	    }
	}
	//fmt.Printf("POINTER :: The address of the routinePool in GetRoutinePool : %p\n", toReturn)
	return toReturn
}

func ShutdownRoutinePools() bool {
	for name, routinePool := range routinePoolMap { 
	    routinePool.log("Shutting down RoutinePool : "+name)
	    close(routinePool.jobChannel)
	    close(routinePool.shutdownChannel)
	}
	return true
}
