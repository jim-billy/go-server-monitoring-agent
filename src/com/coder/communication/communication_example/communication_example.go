package main

import (
	"fmt"
	"com/coder/communication"
	"encoding/json"
)

func testGet(){
	var server communication.Server
	server.Host = "localhost"
	server.Port = 9080
	server.Protocol = communication.HTTP_PROTOCOL
	
	communication.SetDefaultServer(&server)
	connector := communication.GetConnector(communication.HTTP_PROTOCOL)
	request := communication.NewHttpRequest(&server, nil)
	request.Method = communication.HTTP_GET
	request.Api = "/api/monitors"
	response := connector.SendRequest(request)
	fmt.Println("Response ================== ",response)
}

func testPost(){
	var server communication.Server
	server.Host = "localhost"
	server.Port = 9080
	server.Protocol = communication.HTTP_PROTOCOL
	
	communication.SetDefaultServer(&server)
	connector := communication.GetConnector(communication.HTTP_PROTOCOL)
	request := communication.NewHttpRequest(&server, nil)
	request.Method = communication.HTTP_POST
	params, _  := json.Marshal((map[string]string{
		"displayName" : "One",
		"name" : "1",
		"type":"SERVER",
	}))
	request.Data = params
	request.Api = "/api/monitors"
	response := connector.SendRequest(request)
	fmt.Println("Response ================== ",response)
}

func testSendRequest(){
	testGet()	
	testPost()
}

func main(){
	testSendRequest()
}