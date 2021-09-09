
#!/bin/bash 

export SUCCESS=0 
export FAILURE=1 


export WORK_GO_SERVER_MONITORING_AGENT_HOME=`pwd`
export WORK_GO_SERVER_MONITORING_AGENT_BIN=$WORK_GO_SERVER_MONITORING_AGENT_HOME/bin
export WORK_GO_SERVER_MONITORING_AGENT_CMD=$WORK_GO_SERVER_MONITORING_AGENT_HOME/cmd
export WORK_GO_SERVER_MONITORING_AGENT_CMD_MONAGENT=$WORK_GO_SERVER_MONITORING_AGENT_CMD/monagent

export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT=$WORK_GO_SERVER_MONITORING_AGENT_BIN/server-monitoring-agent
export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_CONF=$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/conf
export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_DATA=$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/data
export WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_LOGS=$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/logs
export WORK_GO_SERVER_MONITORING_AGdENT_BIN_SERVER_MONITORING_AGENT_SCRIPTS=$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/scripts


export GO_SERVER_MONITORING_AGENT_BINARY_FILE_NAME='server-monitoring-agent'

#Set src folder in GOPATH
#export GOPATH=$GOPATH:$WORK_GO_SERVER_MONITORING_AGENT_HOME


buildMonAgent() {
	echo "============================== Building Server Monitoring Agent =============================="
	echo "GOROOT : "$GOROOT
	echo "GOPATH : "$GOPATH

	echo "Agent binary output folder : $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT"

	# Check and create binary output folder
	if [ ! -d "$WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT" ]; then
		mkdir -p $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT
	fi

	# Delete old files in the binary output folder
	echo "Deleting old files in the folder $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT"
	# rm -vf $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_CONF/*
	# rm -vf $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_DATA/*
	# rm -vf $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_LOGS/*
	# rm -vf $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_SCRIPTS/*
	
	# Building server monitoring agent
	echo "Building $WORK_GO_SERVER_MONITORING_AGENT_CMD_MONAGENT/monagent.go"
	go build -o $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT/$GO_SERVER_MONITORING_AGENT_BINARY_FILE_NAME $WORK_GO_SERVER_MONITORING_AGENT_CMD_MONAGENT/monagent.go

	if [ "$?" != $SUCCESS ]; then
		echo "*********************** Error while building server monitoring agent **************************"
		exit 1
	fi

	#Copying conf files
	echo "Copying conf files to $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_CONF"
	cp -rf $WORK_GO_SERVER_MONITORING_AGENT_CMD_MONAGENT/conf $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT

	#Copying script files
	echo "Copying script files to $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT_SCRIPTS"
	cp -rf $WORK_GO_SERVER_MONITORING_AGENT_CMD_MONAGENT/scripts $WORK_GO_SERVER_MONITORING_AGENT_BIN_SERVER_MONITORING_AGENT
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




