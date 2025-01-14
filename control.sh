#!/bin/bash



if [ "$1" == "pause" ]; then
    pkill -SIGUSR1 -f "./build/"
    echo "Paused all services."
elif [ "$1" == "resume" ]; then
    pkill -SIGUSR2 -f "./build/"
    echo "Resumed all services."
else
    echo "Usage: $0 {pause|resume}"
fi
