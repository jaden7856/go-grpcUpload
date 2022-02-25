#!/bin/bash

gsTarget=( "AsyngRPCs" )
gsCL=( "cnt" "lop" )
gnPort=50050

pkill -9 go-server

for sTarget in ${gsTarget[*]}
do
	echo ${sTarget}
	cd ${sTarget}

	if [ "$?" == "0" ]
	then
			for nSize in 512 # 1024 2048 4096 8192 16384 32768 65536
			do
				echo "./server/go-server-$sTarget -a=kubespray:$gnPort -size=$nSize &"
				./server/go-server-$sTarget -a=kubespray:$gnPort -size=$nSize &

				if [ "$?" == "0" ]
				then
					echo "OK"	
				else
					echo "Failure"
				fi
			done
	else
		echo "Failure"
	fi

	cd ..
	echo $(pwd)
done
