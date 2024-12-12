#!/usr/bin/env bash

# If workspace isn't clean, it's 'dev'.

if [ "$1" = "main" ]; then
    echo "main"
else
    echo "dev"
fi

