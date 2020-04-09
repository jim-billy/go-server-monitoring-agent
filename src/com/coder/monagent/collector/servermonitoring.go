package collector

import (
	"com/coder/executor"
	"com/coder/initializer"
	"com/coder/monagent/agentconstants"
	"com/coder/routinepool"
)

// ServerMonitoringJob represents the details of the job that has been scheduled
type ServerMonitoringJob struct {
	ID            int
	MonitorConfig LinuxMonitor
	ResultData    *routinepool.JobResult
}

// DoJob implements Job interface's DoJob(routinePool *routinepool.RoutinePool) defined in routinepool.go
// This method will be called in WorkRoutine's safelyDoWork method
func (serverMonJob *ServerMonitoringJob) DoJob(routinePool *routinepool.RoutinePool) {
	linuxmonitor := serverMonJob.MonitorConfig
	jobResult := new(routinepool.JobResult)
	serverMonJob.ResultData = jobResult
	exec := new(executor.Executor)
	agentconstants.Logger.Infof("============================== DoJob : Collecting data : ", linuxmonitor)
	if linuxmonitor.Script {
		agentScriptFilePath := initializer.GetAgentScriptsDir() + "/" + linuxmonitor.Command
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
	GetcollectorAPI().ParseAndSave(serverMonJob)
}

// GetID uniquely identifies the ServerMonitoringJob
func (serverMonJob *ServerMonitoringJob) GetID() int {
	return serverMonJob.ID
}
