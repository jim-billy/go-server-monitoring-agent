package collector

import (
	"com/coder/initializer"
	"com/coder/routinepool"
	"com/coder/executor"
	//"com/coder/logging"
	"fmt"
)

type ServerMonitoringJob struct {
    Id 		string
    MonitorConfig LinuxMonitor
    ResultData *routinepool.JobResult
}

/*
func (serverMonJob ServerMonitoringJob) DoJob(routinePool routinepool.RoutinePool) {
	Logger := logging.GetLogger("url_crawl_agent", "")
    //Logger.Infof("============================== DoJob : Collecting data for Website : %v \n", websiteJob.Website)
    jobResult := new(routinepool.JobResult)
    if(true){
    	//fmt.Println("Data collection using fasthttp")
		jobResult.Result = map[string]interface{}{"urlMetrics": GetUrlMetricsFastHttp(websiteJob)}	
	}else{
		jobResult.Result = map[string]interface{}{"urlMetrics": GetUrlMetrics(websiteJob.Website)}
	}
    websiteJob.ResultData = jobResult
    var intfValue interface{} 
    //urlMet = util.GetValueFromMap(websiteJob.ResultData.Result)
    resultMap := websiteJob.ResultData.Result
    intfValue = resultMap["urlMetrics"]
    //var urlMet util.UrlStats
    //urlMet, ok := intfValue.(*util.UrlStats)
    urlMet := intfValue.(*UrlStats)
    //fmt.Println(urlMet, ok)
//     b, err := json.Marshal(websiteJob.ResultData.Result["urlMetrics"].(interface{}))
//     if err != nil {
//         fmt.Printf("JSON Error: %s", err)
//         return;
//     }
	//fmt.Println("Collected data : Id : %d, website : %s, Status code : %d : ",websiteJob.Id, websiteJob.Website, urlMet.StatusCode)
    Logger.Infof("Collected data : Id : %d, website : %s, Status code : %d : ",websiteJob.Id, websiteJob.Website, urlMet.StatusCode)
    //Logger.Infof("=========================== Collected data === %d === JobId : %d, website : %s, Status code : %d, data : %s : ",totalUrlCounter, websiteJob.Id, urlMet.Url, urlMet.StatusCode, urlMet)
    routinePool.GetCompletedJobsChannel() <- websiteJob
    
}
*/

func (serverMonJob ServerMonitoringJob) DoJob(routinePool *routinepool.RoutinePool) {
	linuxmonitor := serverMonJob.MonitorConfig
	jobResult := new(routinepool.JobResult)
	serverMonJob.ResultData = jobResult
    exec := new(executor.Executor)
    fmt.Println("============================== DoJob : Collecting data : ", linuxmonitor)	
    if(linuxmonitor.Script){
    	agentScriptFilePath := initializer.GetAgentScriptsDir() +"/"+linuxmonitor.Command
	    exec.SetCommand(agentScriptFilePath)	
    }
	exec.SetCommandArgs([]string{linuxmonitor.CommandArgs})
	exec.SetTimeout(10)
	exec.Execute()
	jobResult.Result = make(map[string]interface{})
	serverMonJob.ResultData.Result["is_success"] = exec.IsSuccess()
	serverMonJob.ResultData.Result["execution_time"] = exec.GetExecutionTime()
	serverMonJob.ResultData.Result["error"] = exec.GetError()
	serverMonJob.ResultData.Result["output"] = exec.GetOutput()
	//fmt.Println("===================== Is success : ",exec.IsSuccess(),", Execution time : ",exec.GetExecutionTime(),", Output ",exec.GetOutput(),", Error : ",exec.GetError())
	//fmt.Println("===================== Is success : ",serverMonJob.ResultData.Result["is_success"],", Execution time : ",serverMonJob.ResultData.Result["execution_time"],", Output ",serverMonJob.ResultData.Result["output"],", Error : ",serverMonJob.ResultData.Result["error"])
    GetCollectorApi().ParseAndSave(serverMonJob)
}

