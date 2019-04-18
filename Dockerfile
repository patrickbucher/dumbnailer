FROM golang:1.12.3-stretch AS builder
LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
RUN apt-get update && apt-get install -y git ca-certificates
COPY *.go go.mod /src/
WORKDIR /src
RUN go test && go build -o /app/dumbnailer

FROM debian:stretch-slim
LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
RUN apt-get update && apt-get install -y imagemagick && apt-get autoclean
COPY --from=builder /app/dumbnailer /bin/dumbnailer
ENV DUMBNAILER_PORT="8888" IMAGE_MAGICK="/usr/bin/convert"
EXPOSE $DUMBNAILER_PORT
RUN groupadd -g 1001 gopher && useradd -g 1001 gopher
USER gopher
CMD ["/bin/dumbnailer"]
