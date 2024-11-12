#!/bin/bash

docker build -f Containerfile --no-cache --progress=plain -t $REGISTRY_USERNAME/labr:$(cat version.txt) .
docker push $REGISTRY_USERNAME/labr:$(cat version.txt)
