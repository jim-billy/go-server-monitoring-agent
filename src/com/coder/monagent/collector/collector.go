package collector

import (
	"encoding/json"
	//"fmt"
	"strings"
	//"strconv"
	"com/coder/monagent/agentconstants"
	"com/coder/config"
	"com/coder/scheduler"
	"com/coder/util"
)

var collectorApi CollectorApi
var parserApi ParserApi
var monitorNameVsConfigMap map[string]LinuxMonitor
var metricNameVsParseConfigMap map[string]ParseConfiguration

type CollectorApi struct{
	LinuxMonitors []LinuxMonitor
}

type ParserApi struct{
	
}

type LinuxMonitor struct {
    Name   string `json:"name"`
    Script   bool `json:"script"`
    Command  string    `json:"command"`
    CommandArgs string `json:"command_args"`
    Interval int `json:"interval"`
    KeyValue   bool `json:"key_value"`
    ParseAllLine   bool `json:"parse_all_line"`
    ParseImpl string `json:"parse_impl"`
    ParseConfig []ParseConfiguration `json:"parse_config"`
}

type ParseConfiguration struct { 
	MetricName   string `json:"metric_name"`
    ParseLine   int `json:"parse_line"`
    Delimiter  string    `json:"delimiter"`
    Token int `json:"token"`
    Counter bool `json:"counter"`
    Expression  string  `json:"expression"`
}

type CollectedData struct {
    CollectionTime	int64
    Data   interface{}
    Save   bool
}


func (collectorApi *CollectorApi) Initialize(){
	agentconstants.Logger.Infof("CollectorApi : Initialize : Initializing CollectorApi")
	monitorNameVsConfigMap = make(map[string]LinuxMonitor)
	metricNameVsParseConfigMap = make(map[string]ParseConfiguration)
	collectorApi.loadCollectorConfig()
}

func (collectorApi *CollectorApi) loadCollectorConfig(){
	configLoader := config.GetConfigLoader()
	byteArr, errToReturn := configLoader.LoadBytesFromJson(agentconstants.LINUX_MONITORS_FILE_PATH)
	json.Unmarshal(byteArr, &collectorApi.LinuxMonitors)
	agentconstants.Logger.Infof("CollectorApi : loadCollectorConfig : LinuxMonitors : ", collectorApi.LinuxMonitors)
	if(errToReturn != nil){
		agentconstants.Logger.Infof("CollectorApi : loadCollectorConfig : Error while loading CollectorConfig : ",errToReturn)
	}else{
		for index, _ := range collectorApi.LinuxMonitors {
			linuxMonitor := collectorApi.LinuxMonitors[index]
			agentconstants.Logger.Infof("CollectorApi : loadCollectorConfig : "+linuxMonitor.Name," :::::::: ",linuxMonitor.Interval)
			monitorNameVsConfigMap[linuxMonitor.Name] = linuxMonitor
			parseConfigArr := linuxMonitor.ParseConfig
			for _, parseConfig := range parseConfigArr {
				metricNameVsParseConfigMap[parseConfig.MetricName] = parseConfig
			}
		}
	}
}

func (collectorApi *CollectorApi) ScheduleDataCollection(){
	agentconstants.Logger.Infof("CollectorApi : ScheduleDataCollection : Scheduling data collection")
	var sched *scheduler.Scheduler
	sched = scheduler.GetScheduler("DataCollectionScheduler")
	sched.SetLogger(agentconstants.Logger)
	for index, _ := range collectorApi.LinuxMonitors {
		linuxMonitor := collectorApi.LinuxMonitors[index]
		var schTask scheduler.ScheduleTask
		schTask.SetName(linuxMonitor.Name)
		schTask.SetType(scheduler.REPETITIVE_TASK)
		schTask.SetInterval(linuxMonitor.Interval)
		serverMonJob := &ServerMonitoringJob{MonitorConfig : linuxMonitor, Id : 1}
		schTask.SetJob(serverMonJob)
		sched.Schedule(schTask)
	}
}

func (collectorApi *CollectorApi) ParseAndSave(serverMonJob *ServerMonitoringJob){
	var collectedData *CollectedData
	agentconstants.Logger.Infof("CollectorApi : ParseAndSave : Is success : ",serverMonJob.ResultData.Result["is_success"],", Execution time : ",serverMonJob.ResultData.Result["execution_time"],", Output ",serverMonJob.ResultData.Result["output"],", Error : ",serverMonJob.ResultData.Result["error"])
	collectedData = parserApi.parse(serverMonJob)
	agentconstants.Logger.Infof("CollectorApi : ParseAndSave : Collected data : ",collectedData) 
	//GetFileHandler
}

func (collectorApi *CollectorApi) getParseConfig(name string) (ParseConfiguration){
	return metricNameVsParseConfigMap[name]
}

func (parserApi *ParserApi) parse(serverMonJob *ServerMonitoringJob) *CollectedData{
	var collectedData *CollectedData
	if(serverMonJob.MonitorConfig.KeyValue){
		collectedData = parserApi.parseKeyValue(serverMonJob)
	}else if(serverMonJob.MonitorConfig.ParseAllLine){
		collectedData = parserApi.parseAllLines(serverMonJob)
	}
	return collectedData
}

/*
For the below conf

{
		"name" : "cpu_utilization",
		"script" : true,
		"command" : "metrics.sh",
		"command_args" : "cpu_util",
		"interval" : 5,
		"key_value" : true,
		"parse_impl" : "",
		"parse_config" : [
			{
				"metric_name" : "cpu_instance", 
				"parse_line" : 1,
				"delimiter" : ":",
				"token" : 2,
				"counter" : false 
			},
			{
				"metric_name" : "cpu_idle_percentage", 
				"parse_line" : 2,
				"delimiter" : ":",
				"token" : 2,
				"counter" : false 
			},
			{
				"metric_name" : "cpu_load_percentage", 
				"parse_line" : 3,
				"delimiter" : ":",
				"token" : 2,
				"counter" : false
			},
			{
				"metric_name" : "cpu_wait_percentage", 
				"parse_line" : 4,
				"delimiter" : ":",
				"token" : 2,
				"counter" : false 
			}
		]
			
	},
	
	
Output will be

	{"cpu_idle_percentage ":"84.1","cpu_instance ":"cpu","cpu_load_percentage ":"15.90","cpu_wait_percentage ":"0.00"}
	
*/

func (parserApi *ParserApi) parseKeyValue(serverMonJob *ServerMonitoringJob) *CollectedData{
	parseConfigArr := serverMonJob.MonitorConfig.ParseConfig
	output := serverMonJob.ResultData.Result["output"].(string)
	colData := make(map[string]interface{})
	outputArr := strings.SplitAfter(output, "\n")
	for _, outputLine :=  range outputArr{
		for _, parseConf :=  range parseConfigArr{
			//fmt.Println(parseConf.MetricName," ?????????????? ",outputLine," ?????????????? ",strings.Index(outputLine, parseConf.MetricName))
			if(strings.Index(outputLine, parseConf.MetricName) != -1){
				if(strings.Index(outputLine, parseConf.Delimiter) != -1){
					metricArr := strings.Split(outputLine, parseConf.Delimiter)
					colData[metricArr[0]] = strings.TrimSpace(metricArr[1])
				}
			}
		}
	}
	jsonString, _ := json.Marshal(colData)
	agentconstants.Logger.Infof("CollectorApi : parseKeyValue : Collected JSON  ::::::::::::::::::::: "+serverMonJob.MonitorConfig.Name+" ::::::::::::::: "+string(jsonString))
	//fmt.Println("ParserApi : parseKeyValue : Collected JSON  ::::::::::::::::::::: ",string(jsonString))
	collectedData := NewCollectedData(colData)
	return collectedData
}

/*
	Parse all lines and store them in the metric_name defined in the linux_monitors.json based on delimiter
	For the below parse config
	"parse_config" : [
			{
				"metric_name" : "name", 
				"parse_line" : 1,
				"delimiter" : "==",
				"token" : 2,
				"counter" : false 
			},
			{
				"metric_name" : "size", 
				"parse_line" : 2,
				"delimiter" : "==",
				"token" : 2,
				"counter" : false 
			},
			{
				"metric_name" : "free_space", 
				"parse_line" : 3,
				"delimiter" : "==",
				"token" : 2,
				"counter" : false
			}
		]
		
	Output will be
	"data" : [
			{
				"name" : "/dev", 
				"size" : "1654071296",
				"free_space" : "1654026240"
			},
			{
				"name" : "/run", 
				"size" : "1654071296",
				"free_space" : "1654026240"
			},
			{
				"name" : "/boot", 
				"size" : "1654071296",
				"free_space" : "1654026240"
			}
		]	
	
*/
func (parserApi *ParserApi) parseAllLines(serverMonJob *ServerMonitoringJob) *CollectedData{
	parseConfigArr := serverMonJob.MonitorConfig.ParseConfig
	output := serverMonJob.ResultData.Result["output"].(string)
	outputArr := strings.SplitAfter(output, "\n")
	var colDataArr []map[string]interface{}
	for _, outputLine :=  range outputArr{
		if(strings.Index(outputLine, serverMonJob.MonitorConfig.Name) != -1){
			continue
		}
		colData := make(map[string]interface{})
		//fmt.Println("outputLine ::::::::::::::::::::: ",outputLine)
		//For each output line iterate parseConfigArr, parse values, put them in a map and append it to the colDataArr list	
		for i, parseConf :=  range parseConfigArr{
			if(strings.Index(outputLine, parseConf.Delimiter) != -1){
				metricArr := strings.Split(outputLine, parseConf.Delimiter)
				//fmt.Println(parseConf.MetricName," =============== ",metricArr[i])	
				colData[parseConf.MetricName] = strings.TrimSpace(metricArr[i])
			}
		}
		colDataArr = append(colDataArr, colData)
	}
	collectedData := NewCollectedData(colDataArr)
	jsonString, _ := json.Marshal(colDataArr)
	agentconstants.Logger.Infof("CollectorApi : parseAllLines : Collected JSON  ::::::::::::::::::::: "+serverMonJob.MonitorConfig.Name+" ::::::::::::::: "+string(jsonString))
	//fmt.Println("ParserApi : parseAllLines : Collected JSON  ::::::::::::::::::::: ",serverMonJob.MonitorConfig.Name,":::::::::;",string(jsonString))
	return collectedData
}

func NewCollectedData(data interface{}) *CollectedData{
	collectedData := new(CollectedData)
	collectedData.CollectionTime = util.NowAsUnixMilli()
	collectedData.Data = data
	collectedData.Save = true
	return collectedData
}

func (collectedData *CollectedData) getJsonString() string{
	jsonDataToReturn, _ := json.Marshal(collectedData.Data)
	//fmt.Println("Collected JSON string ::::::::::::::::::::: ",string(jsonDataToReturn))
	return string(jsonDataToReturn)
	
}


func GetCollectorApi() *CollectorApi{
	return &collectorApi
}

