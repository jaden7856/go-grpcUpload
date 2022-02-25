#!/bin/bash

gsTarget=( "AsyngRPCs" )
gsCL=( "cnt" "lop" )
gnPort=50050

pkill -9 go-client

for sTarget in ${gsTarget[*]}
do
	echo ${sTarget}
	cd ${sTarget}

	if [ "$?" == "0" ]
	then
			for nSize in 512 # 512 1024 2048 4096 8192 16384 32768 65536
			do
				for sCL in ${gsCL[*]}
				do
				    if [ "$sCL" == "cnt" ]
				    then
				    	let nCount=1000000	# [vvv] 1000000
				    	let nLoop=1
				    else
				    	let nCount=1
				    	let nLoop=1000000	# [vvv] 1000000
				    fi

					for nIxTest in 1 # 2 3 4 5 # 반복 테스트 횟수 # [vvv] 5
					do
						sLogName=ztime-$sTarget-$nSize.json

						echo "./client/go-client-$sTarget -a=kubespray:$gnPort -size=$nSize -count=$nCount -loop=$nLoop"
						./client/go-client-$sTarget -a=kubespray:$gnPort -size=$nSize -c=$nCount -l=$nLoop

						if [ "$?" == "0" ]
						then
							echo "OK"	
						else
							touch zzz_$sTarget.err
						fi
					done
				done
			done
	fi

	cd ..
done
