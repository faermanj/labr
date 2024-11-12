#!/bin/bash

# Check if an argument is passed
if [ -n "$1" ]; then
    # Check if the argument is a directory
    if [ -d "$1" ]; then
        DIR="$1"
    else
        echo "Error: $1 is not a directory."
        exit 1
    fi
else
    # Use the current working directory as default
    DIR=$(pwd)
fi

echo "Using directory: $DIR"

docker run -it --rm \
    -v "$DIR":/data \
    $REGISTRY_USERNAME/labr:$(cat version.txt) "/data"