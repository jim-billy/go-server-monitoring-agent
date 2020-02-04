package collector

import (
	"com/coder/monagent/agentconstants"
	"com/coder/initializer"
	"com/coder/routinepool"
	"com/coder/executor"
)

type ServerMonitoringJob struct {
    Id 		int
    MonitorConfig LinuxMonitor
    ResultData *routinepool.JobResult
}

func (serverMonJob *ServerMonitoringJob) DoJob(routinePool *routinepool.RoutinePool) {
	linuxmonitor := serverMonJob.MonitorConfig
	jobResult := new(routinepool.JobResult)
	serverMonJob.ResultData = jobResult
    exec := new(executor.Executor)
    agentconstants.Logger.Infof("============================== DoJob : Collecting data : ", linuxmonitor)
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

func (serverMonJob *ServerMonitoringJob) GetId() int{
    return serverMonJob.Id
}

