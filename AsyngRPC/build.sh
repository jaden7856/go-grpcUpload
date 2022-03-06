#!/bin/bash

#ziphome="/home/oracle/prj/bin"

for sTarget in "AsynSocket" "AsyngRPCs"
do
	cd $sTarget
	echo $sTarget

	rm -f *.json *.log	

	if [ "$?" == "0" ]
	then
		cd server
		rm -f ./go-server-$sTarget
		go build -o go-server-$sTarget server.go
#    rm -f ./server
#    go build
#		mv go-server-$sTarget $ziphome
		cd ../client
		rm -f ./go-client-$sTarget
		go build -o go-client-$sTarget client.go
#    rm -f ./client
#    go build
#		mv go-client-$sTarget $ziphome
		cd ..
	fi

	cd ..
	echo `pwd`
done
