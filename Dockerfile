FROM golang:1.12.3-stretch AS builder
LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
RUN apt-get update && apt-get install -y git ca-certificates
COPY *.go go.mod /src/
WORKDIR /src
RUN go build -o /app/dumbnailer

FROM debian:stretch-slim
LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
RUN apt-get update && apt-get install -y imagemagick
COPY --from=builder /app/dumbnailer /bin/dumbnailer
ENV DUMBNAILER_PORT="8888"
ENV IMAGE_MAGICK="/usr/bin/convert"
EXPOSE $DUMBNAILER_PORT
CMD ["/bin/dumbnailer"]
