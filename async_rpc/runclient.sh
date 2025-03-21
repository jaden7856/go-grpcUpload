#!/bin/bash

gsTarget=( "AsyngRPCs" "AsynSocket" )
#gsCL=( "cnt" "lop")
gnPort=50055
nCount=1000000

#  -9 go-client

for sTarget in "${gsTarget[@]}"; do
    echo "Processing $sTarget..."

    if ! cd "$sTarget"; then
        echo "Failed to change directory to $sTarget"
        continue
    fi

    for nSize in 1024; do  # You can extend the sizes here (e.g., 512 1024 2048 ...)
        let gnPort+=1
#       for sCL in "${gsCL[@]}"; do
#        if [ "$sCL" == "cnt" ]; then
#          let nCount=1000000	# [vvv] 1000000
#          let nLoop=1
#        else
#          let nCount=1
#          let nLoop=1000000	# [vvv] 1000000
#        fi

        for nIxTest in 1; do  # 반복 테스트 횟수 (e.g., 1 2 3 4 5 ...)
            sLogName="ztime-$sTarget-$nSize.json"

            echo $(date +%Y%m%d%H%M%S)

            # Execute the client command
            echo "./client/go-client-$sTarget -add=master:$gnPort -size=$nSize -count=$nCount -logtime=$sLogName -debug=0 &"
            ./client/go-client-$sTarget -add=master:$gnPort -size=$nSize -count=$nCount -logtime=$sLogName -debug=0 &

#          echo "./client/client upload-test -add=localhost:$gnPort -d=/home/client-1/file/"
#          ./client/client upload-test -add=localhost:$gnPort -d=/home/client-1/file/

            echo $(date +%Y%m%d%H%M%S)

            # Check if the command was successful
            if [ "$?" -eq 0 ]; then
                echo "OK"
            else
                touch "zzz_$sTarget.err"  # Log the error in a file
                echo "Error occurred during execution"
            fi
        done
    done

    # Return to the parent directory
    cd ..
done