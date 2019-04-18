#!/bin/sh

docker build . -t dumbnailer && \
docker run -p 8888:8888 -it --name dumbnailer --rm dumbnailer
