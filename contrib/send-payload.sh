#!/bin/bash
if [ $# -eq 2 ]; then
    (echo -n payload=; cat $2) | curl -d @- $1
else
    echo Usage: $0 hostname payload
fi
