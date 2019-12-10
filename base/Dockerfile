FROM golang:1.13.5-alpine

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"

RUN apk add imagemagick && \
    addgroup --gid 1001 gophers && \
    adduser -D --uid 1001 -G gophers gopher

USER 1001

ENV SRC_DIR=/home/gopher/src
ENV BIN_DIR=/home/gopher/bin
RUN mkdir $SRC_DIR && mkdir $BIN_DIR

EXPOSE 8888
