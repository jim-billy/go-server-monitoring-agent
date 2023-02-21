package communication

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	//"encoding/json"
	"bytes"
	//"net/url"
	"fmt"
	"reflect"
)

const HTTP_PROTOCOL = "http"
const HTTPS_PROTOCOL = "https"

const HTTP_REQUEST = "HTTP_REQUEST"
const HTTPS_REQUEST = "HTTPS_REQUEST"

const HTTP_GET = "GET"
const HTTP_POST = "POST"
const HTTP_PUT = "PUT"
const HTTP_DELETE = "DELETE"

const SUCCESS_MESSAGE = "Success"
const ERROR_MESSAGE = "Error while sending request"

var httpsConnector HttpsConnector
var httpConnector HttpConnector

var defaultServer *Server
var defaultProxyServer *ProxyServer

//Server
type Server struct {
	Host     string
	Port     int
	Protocol string
}

func (server Server) New(host string, port int, protocol string) Server {
	server.Host = host
	server.Port = port
	server.Protocol = protocol
	return server
}

//Proxy server
type ProxyServer struct {
	Host     string
	Port     int
	Protocol string
	UserName string
	Password string
}

func (proxyServer ProxyServer) New(host string, port int, protocol string, userName string, password string) ProxyServer {
	proxyServer.Host = host
	proxyServer.Port = port
	proxyServer.Protocol = protocol
	proxyServer.UserName = userName
	proxyServer.Password = password
	return proxyServer
}

type Connector interface {
	SetProtocol(protocol string)
	GetProtocol() string
	SetTimeout(timeout int64)
	GetTimeout() int64
	SendRequest(request Request) *Response
}

type BaseConnector struct {
	protocol string
	timeout  int64
}

func (baseConnector *BaseConnector) SetProtocol(protocol string) {
	baseConnector.protocol = protocol
}

func (baseConnector *BaseConnector) GetProtocol() string {
	return baseConnector.protocol
}

func (baseConnector *BaseConnector) SetTimeout(timeout int64) {
	baseConnector.timeout = timeout
}

func (baseConnector *BaseConnector) GetTimeout() int64 {
	return baseConnector.timeout
}

func (baseConnector *BaseConnector) ValidateRequest(req Request) bool {
	var toReturn bool
	if req.GetServer() != nil {
		toReturn = true
	} else {
		fmt.Println("Error while sending request : Unable to send request to server : ", req.GetServer())
	}
	return toReturn
}

type Request interface {
	RequestInit()
	GetServer() *Server
}

type Response struct {
	StatusCode   string
	Error        error
	Message      string
	ResponseData interface{}
}

//HTTP Request

type HttpRequest struct {
	Server          *Server
	ProxyServer     *ProxyServer
	IsProxy         bool
	IsSecure        bool
	Method          string
	Api             string
	Data            []byte
	DataType        string
	ResponseAction  string
	IsParseResponse bool
	Headers         map[string]interface{}
	CustomParams    map[string]interface{}
	Logger          *log.Logger
	LoggerName      string
	UploadFilePath  string
	UploadFileName  string
}

func (httpRequest *HttpRequest) RequestInit() {
	httpRequest.IsSecure = false
}

func (httpRequest *HttpRequest) GetServer() *Server {
	return httpRequest.Server
}

func (httpRequest *HttpRequest) GetProxyServer() *ProxyServer {
	return httpRequest.ProxyServer
}

//HTTPS Request

type HttpsRequest struct {
	HttpRequest
}

func (httpsRequest *HttpsRequest) RequestInit() {
	httpsRequest.IsSecure = true
}

func (httpsRequest *HttpsRequest) GetServer() *Server {
	return httpsRequest.Server
}

func (httpsRequest *HttpsRequest) GetProxyServer() *ProxyServer {
	return httpsRequest.ProxyServer
}

//HttpConnector

type HttpConnector struct {
	BaseConnector
}

func (httpConnector *HttpConnector) getClient() http.Client {
	client := http.Client{
		Timeout: time.Duration(httpConnector.GetTimeout()),
	}
	return client
}

func (httpConnector *HttpConnector) getRequest(requestType string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(requestType, url, body)
}

func (httpConnector *HttpConnector) getTimeout() time.Duration {
	return time.Duration(httpConnector.GetTimeout() * int64(time.Second))
}

func (httpConnector *HttpConnector) setDefaultHeaders(request *http.Request) {
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Accept", "application/json; version=2.0")
}

func (httpConnector *HttpConnector) SendRequest(request Request) *Response {
	var response Response
	if httpConnector.ValidateRequest(request) {
		httpRequest := request.(*HttpRequest)
		fmt.Println("Sending request.....", string(httpRequest.Data))
		urlToGet := httpConnector.getUrl(request)
		fmt.Println("====================== Sending", httpRequest.Method, " request for the url : ", urlToGet, " =========================")
		dataToPost := httpRequest.Data
		req, _ := httpConnector.getRequest(httpRequest.Method, urlToGet, bytes.NewBuffer(dataToPost))
		httpConnector.setDefaultHeaders(req)
		client := httpConnector.getClient()
		res, err := client.Do(req)
		if res != nil {
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			response.StatusCode = res.Status
			fmt.Println("Response Status Code :", res.Status)
			response.ResponseData = string(body)
			fmt.Println("Response : ", string(body))
			response.Message = SUCCESS_MESSAGE
		} else {
			response.Error = err
			fmt.Println("Error : ", err)
		}
	} else {
		response.Message = ERROR_MESSAGE
		fmt.Println("Error while sending request.....")
	}
	return &response
}

func (httpConnector *HttpConnector) getUrl(request Request) string {
	httpRequest := request.(*HttpRequest)
	urlToReturn := fmt.Sprintf("%s://%s:%s%s", httpConnector.protocol, httpRequest.GetServer().Host, strconv.Itoa(httpRequest.GetServer().Port), httpRequest.Api)
	return urlToReturn
}

//HttpsConnector

type HttpsConnector struct {
	BaseConnector
}

func (httpsConnector *HttpsConnector) getClient() http.Client {
	trTls11 := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			MaxVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: trTls11,
		Timeout:   time.Duration(httpsConnector.GetTimeout()),
	}
	return client
}

func (httpsConnector *HttpsConnector) SendRequest(request Request) *Response {
	var response Response
	if httpConnector.ValidateRequest(request) {
		fmt.Println("Sending request.....", request)
	} else {
		fmt.Println("Error while sending request.....")
	}
	return &response
}

//Public methods

func GetConnector(protocol string) Connector {
	if protocol == HTTPS_PROTOCOL {
		httpsConnector.protocol = HTTPS_PROTOCOL
		return &httpsConnector
	} else if protocol == HTTP_PROTOCOL {
		httpConnector.protocol = HTTP_PROTOCOL
		return &httpConnector
	}
	return nil
}

func NewHttpRequest(server *Server, proxyServer *ProxyServer) *HttpRequest {
	httpRequest := &HttpRequest{}
	httpRequest.RequestInit()

	if server != nil {
		httpRequest.Server = server
	}

	return httpRequest
}

func NewHttpsRequest(server *Server, proxyServer *ProxyServer) *HttpsRequest {
	httpsRequest := &HttpsRequest{}
	httpsRequest.RequestInit()
	if server != nil {
		httpsRequest.Server = server
	}
	return httpsRequest
}

func IsInstanceOf(typeName string, typeInterface interface{}) bool {
	fmt.Println("typeInterface : ", reflect.TypeOf(typeInterface), reflect.TypeOf((*HttpsConnector)(nil)))
	switch typeName {
	case HTTP_PROTOCOL:
		return reflect.TypeOf((*HttpConnector)(nil)) == reflect.TypeOf(typeInterface)
	case HTTPS_PROTOCOL:
		return reflect.TypeOf((*HttpsConnector)(nil)) == reflect.TypeOf(typeInterface)
	case HTTP_REQUEST:
		return reflect.TypeOf((*HttpRequest)(nil)) == reflect.TypeOf(typeInterface)
	case HTTPS_REQUEST:
		return reflect.TypeOf((*HttpsRequest)(nil)) == reflect.TypeOf(typeInterface)
	}
	return false
}

func SetDefaultServer(server *Server) {
	defaultServer = server
}

func GetDefaultServer() *Server {
	return defaultServer
}

func SetDefaultProxyServer(proxyServer *ProxyServer) {
	defaultProxyServer = proxyServer
}

func GetDefaultProxyServer() *ProxyServer {
	return defaultProxyServer
}
