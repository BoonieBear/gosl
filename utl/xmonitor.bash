#!/bin/bash

FILES="*.go"

echo
echo "monitoring:"
echo $FILES
echo
echo

refresh(){
    echo
    echo
    echo
    go test -test.run="mylab01"
}

while true; do
    inotifywait -q -e modify $FILES
    refresh
done
