#!/bin/sh

docker build -t dumbnailer-base - <../base/Dockerfile 
docker build -t dumbnailer-s2i .
s2i build https://github.com/patrickbucher/dumbnailer.git dumbnailer-s2i dumbnailer
docker run -p 8888:8888 -dit dumbnailer
