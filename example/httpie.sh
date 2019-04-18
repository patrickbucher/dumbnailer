#!/bin/sh

http -f POST localhost:8888/v1/generatemultiple meta="{ \"page\": 1, \"resolutions\": [ { \"width\": 200, \"height\": 283}, { \"width\": 300, \"height\": 378} ] }" file@document.pdf > response.json && \
    rm -f *.jpg && \
    jq -r '.base64Images[0]' response.json | base64 -d > 1.jpg && \
    jq -r '.base64Images[1]' response.json | base64 -d > 2.jpg && \
    rm response.json
