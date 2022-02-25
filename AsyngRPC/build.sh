#/bin/bash

ziphome="/home/master/go-grpcUpload/bin"

for sTarget in "AsyngRPCs"
do
	cd $sTarget
	echo $sTarget

	rm -f *.json *.log	

	if [ "$?" == "0" ]
	then
		cd server
		rm -f ./go-server-$sTarget
		go build -o go-server-$sTarget server.go
		mv go-server-$sTarget $ziphome
		cd ../client
		rm -f ./go-client-$sTarget
		go build -o go-client-$sTarget client.go
		mv go-client-$sTarget $ziphome
		cd ..
	fi

	cd ..
	echo `pwd`
done
