#!/bin/bash

gsTarget=( "AsynSocket" "AsyngRPCs" )
gsTls=( "tls" "none")
gsCL=( "cnt" "lop" )
gnPort=50056

pkill -9 go-client

for sTarget in ${gsTarget[*]}
do
	echo ${sTarget}
	cd ${sTarget}

	if [ "$?" == "0" ]
	then
		for sTls in ${gsTls[*]}
		do
			if [[ $sTls == "5QUIC" || $sTls == "8AsynQUIC" ]]
			then
				if [ $sTls == "none" ]
				then
					break	# 5QUIC는 1번만 돌면 된다. 기본이 TLS라서
				fi
			fi

			for nSize in 512 # 512 1024 2048 4096 8192 16384 32768 65536
			do
				let gnPort=gnPort+1
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

						echo "./client/go-client-$sTarget -add=master:$gnPort -tls=$sTls -size=$nSize -count=$nCount -loop=$nLoop -logtime=$sLogName -debug=0"
						./client/go-client-$sTarget -add=master:$gnPort -tls=$sTls -size=$nSize -count=$nCount -loop=$nLoop -logtime=$sLogName -debug=0

						if [ "$?" == "0" ]
						then
							echo "OK"
						else
							touch zzz_$sTarget.err
						fi
					done
				done
			done
		done
	fi

	cd ..
done
