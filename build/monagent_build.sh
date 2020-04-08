#!/bin/bash 

export SUCCESS=0 
export FAILURE=1 



export WORK_GO_MY_PUBLIC_PROJECTS=$WORK_GO/go-my-public-projects
export WORK_GO_SERVER_MONITORING_AGENT_HOME=$WORK_GO_MY_PUBLIC_PROJECTS/go-server-monitoring-agent
export WORK_GO_SERVER_MONITORING_AGENT_BIN=$WORK_GO_SERVER_MONITORING_AGENT_HOME/bin

export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT=$WORK_GO_SERVER_MONITORING_AGENT_BIN/server-monitoring-agent
export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_CONF=$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/conf
export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_SCRIPTS=$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/scripts

export WORK_GO_SERVER_MONITORING_AGENT_SRC=$WORK_GO_SERVER_MONITORING_AGENT_HOME/src
export WORK_GO_SERVER_MONITORING_AGENT_SRC_COM_CODER=$WORK_GO_SERVER_MONITORING_AGENT_SRC/com/coder

export GO_SERVER_MONITORING_AGENT_BINARY_FILE_NAME='server-monitoring-agent'

#Set src folder in GOPATH
export GOPATH=$GOPATH:$WORK_GO_SERVER_MONITORING_AGENT_HOME


buildMonAgent() {
	echo "============================== Building Server Monitoring Agent =============================="
	echo "GOROOT : "$GOROOT
	echo "GOPATH : "$GOPATH

	#Building server monitoring agent
	echo "Building $WORK_GO_SERVER_MONITORING_AGENT_SRC_COM_CODER/monagent/monagent.go"
	echo "Agent binary output folder : $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT"
	go build -o $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/$GO_SERVER_MONITORING_AGENT_BINARY_FILE_NAME $WORK_GO_SERVER_MONITORING_AGENT_SRC_COM_CODER/monagent/monagent.go

	if [ "$?" != $SUCCESS ]; then
		echo "*********************** Error while building server monitoring agent **************************"
		exit 1
	fi

	#Copying conf files
	echo "Copying conf files to $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_CONF"
	cp -rf $WORK_GO_SERVER_MONITORING_AGENT_SRC_COM_CODER/monagent/conf $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_CONF

	#Copying script files
	echo "Copying script files to $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_SCRIPTS"
	cp -rf $WORK_GO_SERVER_MONITORING_AGENT_SRC_COM_CODER/monagent/scripts $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_SCRIPTS
}

buildResult() {
	if [ "$?" = $SUCCESS ]; then
		echo ""
		echo "To run the server monitoring agent, please execute the below commands :"
		echo "	cd $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT"
		echo "	./$GO_SERVER_MONITORING_AGENT_BINARY_FILE_NAME"
		echo ""
	else
		echo "*********************** Error while building server monitoring agent **************************"
	fi
}


main() {
	buildMonAgent
	buildResult
}

main

