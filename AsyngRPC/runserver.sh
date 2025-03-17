#!/bin/bash

gsTarget=( "AsyngRPCs" "AsynSocket" )
gnPort=50055

# pkill -9 go-server

# Iterate over targets
for sTarget in "${gsTarget[@]}"; do
    echo "Processing $sTarget..."

    # Change directory to target directory
    if ! cd "$sTarget"; then
        echo "Failed to change directory to $sTarget"
        continue
    fi

    # Iterate over sizes (could be extended to multiple values)
    for nSize in 1024; do
        let gnPort+=1
        echo "Running ./server/go-server-$sTarget -add=master:$gnPort -size=$nSize -debug=0 &"

        # Start the server in the background
        if ./server/go-server-$sTarget -add=master:$gnPort -size=$nSize -debug=0 &; then
            echo "Server started successfully on port $gnPort with size $nSize"
        else
            echo "Failure starting server on port $gnPort with size $nSize"
        fi
    done

    # Return to the previous directory
    cd .. || { echo "Failed to return to previous directory"; exit 1; }

    # Print current working directory for verification
    echo "Returned to $(pwd)"
done