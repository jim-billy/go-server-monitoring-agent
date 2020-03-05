#!/bin/bash 

buildMonAgent() {
	#echo $WORK_GO_BIN
	go build -o $WORK_GO_SRC_COM_CODER/monagent/monagent ./monagent.go
	cp -rf $WORK_GO_SRC_COM_CODER/monagent/conf $WORK_GO_BIN/monagent
	cp -rf $WORK_GO_SRC_COM_CODER/monagent/webagent $WORK_GO_BIN/monagent
}


main() {
	buildMonAgent
}

main

