package url

import (
	"github.com/jim-billy/go-server-monitoring-agent/pkg/logging"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/routinepool"
	//"fmt"
)

type WebsiteJob struct {
	Website    string
	Domain     string
	Id         int
	ResultData *routinepool.JobResult
}

func (websiteJob WebsiteJob) DoJob(routinePool *routinepool.RoutinePool) {
	Logger := logging.GetLogger("url_crawl_agent", "/tmp", true)
	//Logger.Println("============================== DoJob : Collecting data for Website : %v \n", websiteJob.Website)
	jobResult := new(routinepool.JobResult)
	if true {
		//fmt.Println("Data collection using fasthttp")
		jobResult.Result = map[string]interface{}{"urlMetrics": GetUrlMetricsFastHttp(&websiteJob)}
	} else {
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
	Logger.Println("Collected data : Id : %d, website : %s, Status code : %d : ", websiteJob.Id, websiteJob.Website, urlMet.StatusCode)
	//Logger.Println("=========================== Collected data === %d === JobId : %d, website : %s, Status code : %d, data : %s : ",totalUrlCounter, websiteJob.Id, urlMet.Url, urlMet.StatusCode, urlMet)
	routinePool.GetCompletedJobsChannel() <- websiteJob

}

func (websiteJob WebsiteJob) GetID() int {
	return websiteJob.Id
}
