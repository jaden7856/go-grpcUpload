#!/bin/bash

gsTarget=( "AsyngRPCs" "AsynSocket" )
gnPort=50055

# pkill -9 go-server

for sTarget in ${gsTarget[*]}
do
	echo ${sTarget}
	cd ${sTarget}

	if [ "$?" == "0" ]
	then
    for nSize in 1024 # 1024 2048 4096 8192 16384 32768 65536
    do
      let gnPort=gnPort+1
      echo "./server/go-server-$sTarget -add=master:$gnPort -size=$nSize -debug=0 &"
      ./server/go-server-$sTarget -add=master:$gnPort -size=$nSize -debug=0 &

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
	echo `pwd`
done

