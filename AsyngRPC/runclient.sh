#!/bin/bash

gsTarget=( "AsyngRPCs" "AsynSocket" )
#gsCL=( "cnt" "lop")
gnPort=50055

#  -9 go-client

for sTarget in ${gsTarget[*]}
do
	echo ${sTarget}
	cd ${sTarget}

	if [ "$?" == "0" ]
	then

    for nSize in 1024 # 512 1024 2048 4096 8192 16384 32768 65536
    do
      let gnPort=gnPort+1
      let nCount=1000000
#      for sCL in ${gsCL[*]}
#      do
#        if [ "$sCL" == "cnt" ]
#        then
#          let nCount=1000000	# [vvv] 1000000
#          let nLoop=1
#        else
#          let nCount=1
#          let nLoop=1000000	# [vvv] 1000000
#        fi
        for nIxTest in 1 # 2 3 4 5 # 반복 테스트 횟수 # [vvv] 5
        do
          sLogName=ztime-$sTarget-$nSize.json

          echo `date +%Y%m%d%H%M%S`

          echo "./client/go-client-$sTarget -add=master:$gnPort -size=$nSize -count=$nCount -logtime=$sLogName -debug=0 &"
          ./client/go-client-$sTarget -add=master:$gnPort -size=$nSize -count=$nCount -logtime=$sLogName -debug=0 &

#          echo "./client/client upload-test -add=localhost:$gnPort -d=/home/client-1/file/"
#          ./client/client upload-test -add=localhost:$gnPort -d=/home/client-1/file/

          echo `date +%Y%m%d%H%M%S`

          if [ "$?" == "0" ]
          then
            echo "OK"
          else
            touch zzz_$sTarget.err
          fi
        done
      done
	fi

	cd ..
done
