package url

import (
	"fmt"
	"log"

	/*

			"net/http"
		    "net/http/httptrace"
		    "github.com/jim-billy/go-server-monitoring-agent/pkg/shutdown"
		    "time"
		    "crypto/tls"
		    "runtime"

	*/
	"bytes"
	"strconv"
	"strings"

	"github.com/jim-billy/go-server-monitoring-agent/pkg/initializer"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/logging"
	"github.com/jim-billy/go-server-monitoring-agent/pkg/util"

	"github.com/valyala/fasthttp"

	_ "net/http/pprof"
)

var https_success_hosts bytes.Buffer
var https_success_hosts_vs_status_code bytes.Buffer
var error_hosts bytes.Buffer
var error_hosts_vs_error bytes.Buffer

var httpsSuccessHostsFile string
var httpsSuccessHostsVsStatusCodeFile string
var errorHostsFile string
var errorHostsVsErrorFile string

var hosts_200_buffer bytes.Buffer
var hosts_301_buffer bytes.Buffer
var hosts_301_new_location_buffer bytes.Buffer
var hosts_3xx_buffer bytes.Buffer
var hosts_4xx_buffer bytes.Buffer
var hosts_5xx_buffer bytes.Buffer

var hosts_200_file string
var hosts_301_file string
var hosts_301_new_location_file string
var hosts_3xx_file string
var hosts_4xx_file string
var hosts_5xx_file string

var hosts_200_counter int
var hosts_301_counter int
var hosts_3xx_counter int
var hosts_4xx_counter int
var hosts_5xx_counter int

var successUrlCounter int
var errorUrlCounter int
var totalUrlCounter int

var ErrorUrlLogger *log.Logger
var SuccessUrlLogger *log.Logger
var StatsLogger *log.Logger

func init() {
	fmt.Println("=================== url_response_util ============== ")
	hosts_200_buffer = bytes.Buffer{}
	hosts_301_buffer = bytes.Buffer{}
	hosts_301_new_location_buffer = bytes.Buffer{}
	hosts_3xx_buffer = bytes.Buffer{}
	hosts_4xx_buffer = bytes.Buffer{}
	hosts_5xx_buffer = bytes.Buffer{}
	https_success_hosts = bytes.Buffer{}
	https_success_hosts_vs_status_code = bytes.Buffer{}
	error_hosts = bytes.Buffer{}
	error_hosts_vs_error = bytes.Buffer{}
}

func WriteDataCollectionStatsToFile() {
	hosts_200_file = initializer.GetAgentDataDir() + "/" + "200_hosts.csv"
	hosts_301_file = initializer.GetAgentDataDir() + "/" + "301_hosts.csv"
	hosts_301_new_location_file = initializer.GetAgentDataDir() + "/" + "301_new_location_hosts.csv"
	hosts_3xx_file = initializer.GetAgentDataDir() + "/" + "3xx_hosts.csv"
	hosts_4xx_file = initializer.GetAgentDataDir() + "/" + "4xx_hosts.csv"
	hosts_5xx_file = initializer.GetAgentDataDir() + "/" + "5xx_hosts.csv"
	httpsSuccessHostsFile = initializer.GetAgentDataDir() + "/" + "https_hosts.csv"
	httpsSuccessHostsVsStatusCodeFile = initializer.GetAgentDataDir() + "/" + "https_hosts_vs_status_code.csv"
	errorHostsFile = initializer.GetAgentDataDir() + "/" + "error_hosts.csv"
	errorHostsVsErrorFile = initializer.GetAgentDataDir() + "/" + "error_hosts_vs_error.csv"
	fmt.Println("====================== WriteDataCollectionStatsToFile ======================")
	//util.AppendToFile(httpsSuccessHostsFile, https_success_hosts.String())

	err := util.WriteToFile(httpsSuccessHostsFile, https_success_hosts.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", httpsSuccessHostsFile, err)
	}
	err = util.WriteToFile(httpsSuccessHostsVsStatusCodeFile, https_success_hosts_vs_status_code.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", httpsSuccessHostsVsStatusCodeFile, err)
	}
	err = util.WriteToFile(errorHostsFile, error_hosts.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", errorHostsFile, err)
	}
	err = util.WriteToFile(errorHostsVsErrorFile, error_hosts_vs_error.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", errorHostsVsErrorFile, err)
	}

	err = util.WriteToFile(hosts_200_file, hosts_200_buffer.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", hosts_301_file, err)
	}

	err = util.WriteToFile(hosts_301_file, hosts_301_buffer.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", hosts_301_file, err)
	}

	err = util.WriteToFile(hosts_301_new_location_file, hosts_301_new_location_buffer.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", hosts_301_new_location_file, err)
	}

	err = util.WriteToFile(hosts_3xx_file, hosts_3xx_buffer.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", hosts_3xx_file, err)
	}

	err = util.WriteToFile(hosts_4xx_file, hosts_4xx_buffer.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", hosts_4xx_file, err)
	}

	err = util.WriteToFile(hosts_5xx_file, hosts_5xx_buffer.String())
	if err != nil {
		fmt.Println("Error while wrting to file : ", hosts_5xx_file, err)
	}

	StatsLogger.Println("==================================== Data collection stats ====================================")
	StatsLogger.Println("hosts_200_counter : " + strconv.Itoa(hosts_200_counter))
	StatsLogger.Println("hosts_301_counter : " + strconv.Itoa(hosts_301_counter))
	StatsLogger.Println("hosts_3xx_counter : " + strconv.Itoa(hosts_3xx_counter))
	StatsLogger.Println("hosts_4xx_counter : " + strconv.Itoa(hosts_4xx_counter))
	StatsLogger.Println("hosts_5xx_counter : " + strconv.Itoa(hosts_5xx_counter))
	StatsLogger.Println("successUrlCounter : " + strconv.Itoa(successUrlCounter))
	StatsLogger.Println("errorUrlCounter : " + strconv.Itoa(errorUrlCounter))
	StatsLogger.Println("totalUrlCounter : " + strconv.Itoa(totalUrlCounter))

}

//func ParseFastHttpResponse(err error, resp *fasthttp.Response){
func ParseFastHttpResponse(err error, resp *fasthttp.Response, urlStats *UrlStats, websiteJob *WebsiteJob) {
	//fmt.Println("error : ",err,resp)
	domain := websiteJob.Domain
	index := websiteJob.Id
	url := websiteJob.Website
	if err == nil {
		successUrlCounter += 1
		if resp.StatusCode() == 200 {
			hosts_200_buffer.WriteString(domain)
			hosts_200_buffer.WriteString(",")
			hosts_200_counter += 1
		} else if resp.StatusCode() == 301 || resp.StatusCode() == 302 {
			hosts_301_buffer.WriteString(domain)
			hosts_301_buffer.WriteString(",")
			hosts_301_counter += 1
			new_3xx_host := strings.Replace(string(resp.Header.Peek("Location")), "https://", "", -1)
			hosts_301_new_location_buffer.WriteString(new_3xx_host)
			hosts_301_new_location_buffer.WriteString(",")
		} else if resp.StatusCode() > 302 && resp.StatusCode() <= 399 {
			hosts_3xx_buffer.WriteString(domain)
			hosts_3xx_buffer.WriteString(",")
			hosts_3xx_counter += 1
		} else if resp.StatusCode() >= 400 && resp.StatusCode() <= 499 {
			hosts_4xx_buffer.WriteString(domain)
			hosts_4xx_buffer.WriteString(",")
			hosts_4xx_counter += 1
		} else if resp.StatusCode() >= 500 && resp.StatusCode() <= 599 {
			hosts_5xx_buffer.WriteString(domain)
			hosts_5xx_buffer.WriteString(",")
			hosts_5xx_counter += 1
		}

		https_success_hosts.WriteString(domain)
		https_success_hosts.WriteString(",")
		https_success_hosts_vs_status_code.WriteString(strconv.Itoa(index))
		https_success_hosts_vs_status_code.WriteString("***********")
		https_success_hosts_vs_status_code.WriteString(url)
		https_success_hosts_vs_status_code.WriteString("***********")
		https_success_hosts_vs_status_code.WriteString(strconv.Itoa(resp.StatusCode()))
		https_success_hosts_vs_status_code.WriteString("\n")
		//fmt.Println("Id : ",index,", Url : ",url,", Status code : ",resp.StatusCode(),", urlStats : ",urlStats)
		SuccessUrlLogger.Println("Id : ", index, ", Url : ", url, ", Status code : ", resp.StatusCode(), ", urlStats : ", urlStats)
	} else {
		errorUrlCounter += 1
		error_hosts.WriteString(domain)
		error_hosts.WriteString(",")
		error_hosts_vs_error.WriteString(strconv.Itoa(index))
		error_hosts_vs_error.WriteString("***********")
		error_hosts_vs_error.WriteString(url)
		error_hosts_vs_error.WriteString("***********")
		error_hosts_vs_error.WriteString(err.Error())
		error_hosts_vs_error.WriteString("\n")
		ErrorUrlLogger.Println("Id : ", index, ", Url : ", url, ", err : ", err)
		//fmt.Println("Id : ",index,", Url : ",url,", err : ",err)
	}

	totalUrlCounter += 1

}

func SetLogging() {
	StatsLogger = logging.GetLogger("url_collection_stats", "/tmp", true)
	SuccessUrlLogger = logging.GetLogger("url_success", "/tmp", true)
	ErrorUrlLogger = logging.GetLogger("url_error", "/tmp", true)
}
