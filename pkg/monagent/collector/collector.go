package collector

import (
	"encoding/json"
	"fmt"
	"strconv"

	//"fmt"
	"strings"
	//"strconv"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/communication"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/config"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/initializer"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/monagent/agentconstants"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/scheduler"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/util"
)

var collectorAPI CollectorAPI
var parserAPI ParserAPI
var monitorNameVsConfigMap map[string]LinuxMonitor
var metricNameVsParseConfigMap map[string]ParseConfiguration

// CollectorAPI is responsible for data collection in the server
type CollectorAPI struct {
	LinuxMonitors []LinuxMonitor
}

// ParserAPI is responsible for parsing the collected data
type ParserAPI struct {
}

// LinuxMonitor holds all the attributes related to the data collection of metrics like CPU utilization, disk utilization, memory utilization etc.
type LinuxMonitor struct {
	Name         string               `json:"name"`
	Script       bool                 `json:"script"`
	Command      string               `json:"command"`
	CommandArgs  string               `json:"command_args"`
	Interval     int                  `json:"interval"`
	KeyValue     bool                 `json:"key_value"`
	ParseAllLine bool                 `json:"parse_all_line"`
	ParseImpl    string               `json:"parse_impl"`
	ParseConfig  []ParseConfiguration `json:"parse_config"`
}

// ParseConfiguration holds the attributes defined in the linux_monitors.json that are necessary for parsing the collected data.
type ParseConfiguration struct {
	MetricName string `json:"metric_name"`
	ParseLine  int    `json:"parse_line"`
	Delimiter  string `json:"delimiter"`
	Token      int    `json:"token"`
	Counter    bool   `json:"counter"`
	Expression string `json:"expression"`
}

// CollectedData holds the data that is collected
type CollectedData struct {
	Name           string      `json:"name"`
	CollectionTime int64       `json:"collection_time"`
	Data           interface{} `json:"data"`
	Save           bool        `json:"-"` //Ignores this attribute while marshalling
}

// Initialize is responsible loading the data collection metrics from the linux_monitors.json and initialzing the collectorAPI
func (collectorAPI *CollectorAPI) Initialize() {
	agentconstants.Logger.Println("collectorAPI : Initialize : Initializing collectorAPI")
	monitorNameVsConfigMap = make(map[string]LinuxMonitor)
	metricNameVsParseConfigMap = make(map[string]ParseConfiguration)
	collectorAPI.loadCollectorConfig()
}

func (collectorAPI *CollectorAPI) loadCollectorConfig() {
	configLoader := config.GetConfigLoader()
	byteArr, errToReturn := configLoader.LoadBytesFromJson(agentconstants.LinuxMonitorsFilePath)
	json.Unmarshal(byteArr, &collectorAPI.LinuxMonitors)
	agentconstants.Logger.Println("collectorAPI : loadCollectorConfig : LinuxMonitors : ", collectorAPI.LinuxMonitors)
	if errToReturn != nil {
		agentconstants.Logger.Println("collectorAPI : loadCollectorConfig : Error while loading CollectorConfig : ", errToReturn)
	} else {
		for index := range collectorAPI.LinuxMonitors {
			linuxMonitor := collectorAPI.LinuxMonitors[index]
			agentconstants.Logger.Println("collectorAPI : loadCollectorConfig : "+linuxMonitor.Name, " :::::::: ", linuxMonitor.Interval)
			monitorNameVsConfigMap[linuxMonitor.Name] = linuxMonitor
			parseConfigArr := linuxMonitor.ParseConfig
			for _, parseConfig := range parseConfigArr {
				metricNameVsParseConfigMap[parseConfig.MetricName] = parseConfig
			}
		}
	}
}

// ScheduleDataCollection is responsible for scheduling the metrics defined in linux_monitors.json based on the time interval configured.
func (collectorAPI *CollectorAPI) ScheduleDataCollection() {
	agentconstants.Logger.Println("collectorAPI : ScheduleDataCollection : Scheduling data collection")
	var sched *scheduler.Scheduler
	sched = scheduler.GetScheduler("DataCollectionScheduler")
	sched.SetLogger(agentconstants.Logger)
	for index := range collectorAPI.LinuxMonitors {
		linuxMonitor := collectorAPI.LinuxMonitors[index]
		var schTask scheduler.ScheduleTask
		schTask.SetName(linuxMonitor.Name)
		schTask.SetType(scheduler.REPETITIVE_TASK)
		schTask.SetInterval(linuxMonitor.Interval)
		serverMonJob := &ServerMonitoringJob{MonitorConfig: linuxMonitor, ID: 1}
		schTask.SetJob(serverMonJob)
		sched.Schedule(schTask)
	}
}

//ParseAndSave is responsible for parsing and saving the collected data
func (collectorAPI *CollectorAPI) ParseAndSave(serverMonJob *ServerMonitoringJob) {
	var collectedData *CollectedData
	agentconstants.Logger.Println("collectorAPI : ParseAndSave : Is success : ", serverMonJob.ResultData.Result["is_success"], ", Execution time : ", serverMonJob.ResultData.Result["execution_time"], ", Output ", serverMonJob.ResultData.Result["output"], ", Error : ", serverMonJob.ResultData.Result["error"])
	collectedData = parserAPI.parse(serverMonJob)
	agentconstants.Logger.Println("collectorAPI : ParseAndSave : Collected data : ", collectedData)
	//GetFileHandler
}

func (collectorAPI *CollectorAPI) getParseConfig(name string) ParseConfiguration {
	return metricNameVsParseConfigMap[name]
}

func (parserAPI *ParserAPI) parse(serverMonJob *ServerMonitoringJob) *CollectedData {
	var collectedData *CollectedData
	if serverMonJob.MonitorConfig.KeyValue {
		collectedData = parserAPI.parseKeyValue(serverMonJob)
	} else if serverMonJob.MonitorConfig.ParseAllLine {
		collectedData = parserAPI.parseAllLines(serverMonJob)
	}
	collectedData.save()
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

	{
	"name": "cpu_utilization",
	"collection_time": 1607271135099,
	"data": {
			"cpu_idle_percentage ": "96",
			"cpu_instance ": "cpu",
			"cpu_load_percentage ": "4.00",
			"cpu_wait_percentage ": "0.00"
		}
	}

*/

func (parserAPI *ParserAPI) parseKeyValue(serverMonJob *ServerMonitoringJob) *CollectedData {
	parseConfigArr := serverMonJob.MonitorConfig.ParseConfig
	output := serverMonJob.ResultData.Result["output"].(string)
	colData := make(map[string]interface{})
	outputArr := strings.SplitAfter(output, "\n")
	for _, outputLine := range outputArr {
		for _, parseConf := range parseConfigArr {
			//fmt.Println(parseConf.MetricName," ?????????????? ",outputLine," ?????????????? ",strings.Index(outputLine, parseConf.MetricName))
			if strings.Index(outputLine, parseConf.MetricName) != -1 {
				if strings.Index(outputLine, parseConf.Delimiter) != -1 {
					metricArr := strings.Split(outputLine, parseConf.Delimiter)
					colData[strings.TrimSpace(metricArr[0])] = strings.TrimSpace(metricArr[1])
				}
			}
		}
	}
	collectedData := NewCollectedData(serverMonJob.MonitorConfig.Name, colData)
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
	{
		"name": "disk_utilization",
		"collection_time": 1607271133919,
		"data": [
		{
			"free_space": "8252485632",
			"name": "/dev",
			"size": "8252485632"
			},
			{
			"free_space": "7807504384",
			"name": "/dev/shm",
			"size": "8270274560"
			},
			{
			"free_space": "8259801088",
			"name": "/run",
			"size": "8270274560"
			},
			{
			"free_space": "8270274560",
			"name": "/sys/fs/cgroup",
			"size": "8270274560"
			},
			{
			"free_space": "25641586688",
			"name": "/",
			"size": "53660876800"
			},
			{
			"free_space": "760909824",
			"name": "/boot",
			"size": "1063256064"
			},
			{
			"free_space": "384930234368",
			"name": "/home",
			"size": "436736655360"
			},
			{
			"free_space": "1653997568",
			"name": "/run/user/1000",
			"size": "1654054912"
			},
			{}
		]
	}

*/
func (parserAPI *ParserAPI) parseAllLines(serverMonJob *ServerMonitoringJob) *CollectedData {
	parseConfigArr := serverMonJob.MonitorConfig.ParseConfig
	output := serverMonJob.ResultData.Result["output"].(string)
	outputArr := strings.SplitAfter(output, "\n")
	var colDataArr []map[string]interface{}
	for _, outputLine := range outputArr {
		if strings.Index(outputLine, serverMonJob.MonitorConfig.Name) != -1 {
			continue
		}
		colData := make(map[string]interface{})
		//fmt.Println("outputLine ::::::::::::::::::::: ", outputLine)
		//For each output line iterate parseConfigArr, parse values, put them in a map and append it to the colDataArr list
		for i, parseConf := range parseConfigArr {
			if strings.Index(outputLine, parseConf.Delimiter) != -1 {
				metricArr := strings.Split(outputLine, parseConf.Delimiter)
				//fmt.Println(parseConf.MetricName, " =============== ", metricArr[i])
				if parseConf.MetricName != "" {
					colData[parseConf.MetricName] = strings.TrimSpace(metricArr[i])
				}
			}
		}
		if strings.TrimSpace(outputLine) != "" {
			colDataArr = append(colDataArr, colData)
		}

	}
	collectedData := NewCollectedData(serverMonJob.MonitorConfig.Name, colDataArr)
	return collectedData
}

// NewCollectedData returns an instance of the CollectedData struct
func NewCollectedData(name string, data interface{}) *CollectedData {
	collectedData := new(CollectedData)
	collectedData.Name = name
	collectedData.CollectionTime = util.NowAsUnixMilli()
	collectedData.Data = data
	collectedData.Save = true
	return collectedData
}

/*
	Returns the full path of the data collection file to be save
*/
func (collectedData *CollectedData) getDataCollectionFileName() string {
	fileName := initializer.GetAgentDataDir() + "/" + "data_" + collectedData.Name + "_" + strconv.FormatInt(collectedData.CollectionTime, 10)
	return fileName

}

/*
	Saves the collected data to the specified file in the data directory
*/

func (collectedData *CollectedData) save() bool {

	byteJson, err := json.Marshal(collectedData)
	if err != nil {
		agentconstants.Logger.Println("collectorAPI : Error while converting collected data to JSON : " + collectedData.Name + " ::::::::::::::: " + err.Error())
	}
	agentconstants.Logger.Println("collectorAPI : parseKeyValue : Collected JSON  ::::::::::::::::::::: " + collectedData.Name + " ::::::::::::::: " + string(byteJson))
	util.WriteToFile(collectedData.getDataCollectionFileName(), string(byteJson))

	switch collectedData.Name {
	case "mem_utilization":
		postData(byteJson, "/server/metrics/memory")
	case "cpu_utilization":
		postData(byteJson, "/server/metrics/cpu")
	case "disk_utilization":
		postData(byteJson, "/server/metrics/disk")
	}
	//testPostCPU()
	//testPostDisk()
	return true

}

func postData(byteJson []byte, api string) {
	var server communication.Server
	server.Host = "localhost"
	server.Port = 8888
	server.Protocol = communication.HTTP_PROTOCOL

	communication.SetDefaultServer(&server)
	connector := communication.GetConnector(communication.HTTP_PROTOCOL)
	request := communication.NewHttpRequest(&server, nil)
	request.Method = communication.HTTP_POST

	request.Data = byteJson
	request.Api = api
	response := connector.SendRequest(request)
	fmt.Println("Response ================== ", response)
}

func testPostCPU() {
	var server communication.Server
	server.Host = "localhost"
	server.Port = 8888
	server.Protocol = communication.HTTP_PROTOCOL

	communication.SetDefaultServer(&server)
	connector := communication.GetConnector(communication.HTTP_PROTOCOL)
	request := communication.NewHttpRequest(&server, nil)
	request.Method = communication.HTTP_POST
	params, _ := json.Marshal((map[string]string{
		"name":                "cpu",
		"cpu_instance":        "cpu0",
		"cpu_idle_percentage": "22",
		"cpu_load_percentage": "50",
		"cpu_wait_percentage": "15",
	}))
	request.Data = params
	request.Api = "/server/metrics/cpu"
	response := connector.SendRequest(request)
	fmt.Println("Response ================== ", response)
}

func testPostDisk() {
	var server communication.Server
	server.Host = "localhost"
	server.Port = 8888
	server.Protocol = communication.HTTP_PROTOCOL

	communication.SetDefaultServer(&server)
	connector := communication.GetConnector(communication.HTTP_PROTOCOL)
	request := communication.NewHttpRequest(&server, nil)
	request.Method = communication.HTTP_POST
	/*
		params, _ := json.Marshal((map[string]string{
			"name":      "disk",
			"disk_list": "",
		}))
	*/
	params := []byte(`{
		"name": "disk",
		"disk_list": [
			{
				"disk_name": "/opt",
				"free_space": "123123",
				"size": "333333"
			},
			{
				"disk_name": "/run",
				"free_space": "444444",
				"size": "111111"
			},
			{
				"disk_name": "/home",
				"free_space": "888888",
				"size": "7777777"
			}
		]
	  }`)

	request.Data = params
	request.Api = "/server/metrics/disk"
	response := connector.SendRequest(request)
	fmt.Println("Response ================== ", response)
}

// GetcollectorAPI can be invoked from other packages for fetching an instance of the collectorAPI
func GetcollectorAPI() *CollectorAPI {
	return &collectorAPI
}
