package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/logging"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/routinepool"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/shutdown"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/url"

	//"strconv"
	"strings"
	"time"
	//"encoding/json"
	//"math/rand"
)

const (
	ROUTINE_POOL_SIZE = 5
	QUEUE_CAPACITY    = 100
	JOB_THRESHOLD     = 20
)

var Logger *log.Logger

func init() {
	fmt.Println("Init of routinepooltest")
	runtime.GOMAXPROCS(3 * runtime.NumCPU())
	url.LoadWebsites("")
	Logger = logging.GetLogger("agent", "/tmp", true) //Third boolean parameter is for deleting existing logs
	shutdown.GetShutdownHandler().Init(nil)
}

type urlStatus struct {
	url      string
	status   bool
	urlstats *url.UrlStats
}

type WebsiteJob struct {
	Website    string
	Id         int
	ResultData *routinepool.JobResult
}

func (websiteJob *WebsiteJob) DoJob(routinePool *routinepool.RoutinePool) {
	fmt.Println("============================== DoJob : Collecting data for Website : %v \n", websiteJob.Website)
	//time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
	time.Sleep(1 * time.Second)
	/*
		    var a[3] int
			j := 5
			fmt.Println(a[j])
	*/
	/*
	   jobResult := new(routinepool.JobResult)
	   jobResult.Result = map[string]interface{}{"urlMetrics": url.GetUrlMetrics(websiteJob.Website)}
	   websiteJob.ResultData = jobResult
	   var intfValue interface{}
	   resultMap := websiteJob.ResultData.Result
	   intfValue = resultMap["urlMetrics"]

	   urlMet := intfValue.(*url.UrlStats)

	   fmt.Println("=========================== Collected data : Id : ",websiteJob.Id," website : ",urlMet.Url," Status code : ",urlMet.StatusCode," data : ", urlMet)
	*/
	//routinePool.GetCompletedJobsChannel() <- websiteJob

}

func (websiteJob *WebsiteJob) GetID() int {
	return websiteJob.Id
}

func main() {
	//url.SetLimit()
	routinePoolTest()
	shutdown.Wait()
}

func routinePoolTest() {
	config := routinepool.RoutinePoolConfig{
		RoutinePoolName: "DataCollectionPool",
		RoutinePoolSize: ROUTINE_POOL_SIZE,
		QueueCapacity:   QUEUE_CAPACITY,
		//Logger: Logger,
	}
	routinePool, err := routinepool.NewRoutinePool(config)
	if err != nil {
		fmt.Println("Error while creating routinepool : ", err)
	}

	fmt.Printf("POINTER :: The address of the received routinePool in routinePoolTest : %p\n", routinePool)
	//go printRoutinePoolStats()
	time.Sleep(1 * time.Second)
	go sendWebsiteJobs()
	//go sendShutdownSignal()

}

func sendWebsiteJobs() {
	routinePool := routinepool.GetRoutinePool("DataCollectionPool")
	urls := url.GetWebsiteList()
	for i, url := range urls {
		//fmt.Println("%d : %s ",i,url)
		var urlStr string
		if strings.HasPrefix(url, "www.") {
			urlStr = "https://" + url
		} else {
			urlStr = "https://www." + url
		}
		urlStr = "https://127.0.0.1/index.html"
		websiteJob := &WebsiteJob{Website: urlStr, Id: i}
		//logging.Logger.Println("============================ Sending Job : %d %s \n", i,url)
		fmt.Println("============================ Sending Job : ", i, url)
		//fmt.Printf("POINTER :: The address of the received routinePool in sendWebsiteJobs : %p\n", routinePool)
		routinePool.ExecuteJob(websiteJob)
		//time.Sleep(1 * time.Second)
		if i == JOB_THRESHOLD {
			fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% STOPPED SENDING JOBS : JOB_THRESHOLD ", JOB_THRESHOLD, "REACHED %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
			break
		}
	}
}

func printRoutinePoolStats() {
	routinePool := routinepool.GetRoutinePool("DataCollectionPool")
	for {
		routinePool.PerformanceStats()
		time.Sleep(1 * time.Second)
	}
}

func sendShutdownSignal() {
	routinePool := routinepool.GetRoutinePool("DataCollectionPool")
	time.Sleep(5 * time.Second)
	routinePool.GetShutdownChannel()
	//close(routinePool.GetShutdownChannel())
	routinepool.ShutdownRoutinePools()
	//logging.Logger.Println("routinepool.ShutdownRoutinePools() : ",)
	time.Sleep(3 * time.Second)
	//logging.Logger.Println("================== Sending jobs after shutting down the channel after 3 seconds ====================")
	//go sendJobs(routinePool)
}
