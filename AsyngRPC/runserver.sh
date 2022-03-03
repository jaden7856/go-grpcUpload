#!/bin/bash

gsTarget=( "AsynSocket" "AsyngRPCs" )
gsTls=( "tls" "none")
gsCL=( "cnt" "lop" )
gnPort=50050

pkill -9 go-server

for sTarget in ${gsTarget[*]}
do
	echo ${sTarget}
	cd ${sTarget}

	if [ "$?" == "0" ]
	then
		for sTls in ${gsTls[*]}
		do
			if [[ $sTarget == "5QUIC" || $sTarget == "8AsynQUIC" ]]
			then
				if [ $sTls == "none" ]
				then
					break	# 5QUIC는 1번만 돌면 된다. 기본이 TLS라서
				fi
			fi

			for nSize in 512 # 1024 2048 4096 8192 16384 32768 65536
			do
				let gnPort=gnPort+1
				echo "./server/go-server-$sTarget -add=master:$gnPort -size=$nSize -tls=$sTls -debug=0 &"
				./server/go-server-$sTarget -add=master:$gnPort -size=$nSize -tls=$sTls -debug=0 &

				if [ "$?" == "0" ]
				then
					echo "OK"	
				else
					echo "Failure"
				fi
			done
		done
	else
		echo "Failure"
	fi

	cd ..
	echo `pwd`
done
