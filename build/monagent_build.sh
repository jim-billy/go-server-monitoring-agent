#!/bin/bash 

buildMonAgent() {
	#echo $WORK_GO_BIN
	go build -o $MONITORING_SRC_COM_CODER/monagent/monagent ./monagent.go
	cp -rf $MONITORING_SRC_COM_CODER/monagent/conf $WORK_GO_BIN/monagent
	cp -rf $MONITORING_SRC_COM_CODER/monagent/webagent $WORK_GO_BIN/monagent
}


main() {
	buildMonAgent
}

main

