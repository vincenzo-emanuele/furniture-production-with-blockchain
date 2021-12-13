#!/bin/bash


export FABRIC_CFG_PATH=$PWD/../config

function start(){
	echo Avvio della rete
	./network.sh up
	./network.sh createChannel12 -c chanorg1org2
	./network.sh createChannel23 -c chanorg2org3
}

if [[ $# == 0 ]] ; then
	start
elif [[ $1 == "deploy" ]] ; then
	start
	echo Comincio con il deploy
	./network.sh deployCC -c chanorg1org2 -ccn basic-12 -ccp ../asset-transfer-basic/chaincode-go-12 -ccl go
	./network.sh deployCC23 -c chanorg2org3 -ccn basic-23 -ccp ../asset-transfer-basic/chaincode-go-23 -ccl go
fi
