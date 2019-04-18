#!/bin/sh

http -f POST localhost:8888/v1/generatemultiple meta="{ \"page\": 1, \"resolutions\": [ { \"width\": 100, \"height\": 141 }, { \"width\": 200, \"height\": 242 } ] }" file@document.pdf > response.json
