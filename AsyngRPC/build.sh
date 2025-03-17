#!/bin/bash

#ziphome="/home/oracle/prj/bin"

# Terminate existing processes
pkill -9 go-server
pkill -9 go-client

for sTarget in "AsynSocket" "AsyngRPCs"; do
    echo "Processing $sTarget..."

    # Change to the target directory and handle errors
    if ! cd "$sTarget"; then
        echo "Failed to change directory to $sTarget"
        continue
    fi

    # Clean up existing files
    echo "Cleaning up old files..."
    rm -f *.json *.log
    if [ $? -ne 0 ]; then
        echo "Error cleaning up files in $sTarget"
        continue
    fi

    # Compile the server and client
    if cd server && rm -f ./go-server-$sTarget && go build -o go-server-$sTarget server.go; then
        echo "Server compiled successfully for $sTarget"
    else
        echo "Failed to compile server for $sTarget"
        cd ..
        continue
    fi

    if cd ../client && rm -f ./go-client-$sTarget && go build -o go-client-$sTarget client.go; then
        echo "Client compiled successfully for $sTarget"
    else
        echo "Failed to compile client for $sTarget"
        cd ..
        continue
    fi

    # Return to the previous directory
    cd ..
    echo "Returned to $(pwd)"
done