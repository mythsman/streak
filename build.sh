#!/usr/bin/env bash

docker buildx build -t mythsman/streak:amd64-latest -f  Dockerfile --platform=linux/amd64 .

docker buildx build -t mythsman/streak:arm64-latest -f  Dockerfile --platform=linux/arm64 .

docker manifest annotate mythsman/streak  mythsman/streak:amd64-latest --os linux --arch amd64

docker manifest annotate mythsman/streak  mythsman/streak:arm64-latest --os linux --arch arm64 --variant v8

docker push mythsman/streak:amd64-latest

docker push mythsman/streak:arm64-latest

sleep 1

docker manifest rm mythsman/streak

sleep 1

docker manifest create mythsman/streak mythsman/streak:arm64-latest mythsman/streak:amd64-latest

sleep 1

docker manifest push mythsman/streak
