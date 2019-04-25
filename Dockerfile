FROM golang:1.12.4-alpine AS builder
LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
RUN apk add alpine-sdk
COPY *.go go.mod /src/
WORKDIR /src
RUN go test && go build -o /app/dumbnailer

FROM alpine:latest
LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
RUN apk add imagemagick
COPY --from=builder /app/dumbnailer /bin/dumbnailer
ENV PORT="8888" IMAGE_MAGICK="/usr/bin/convert"
EXPOSE $PORT
RUN addgroup -g 1001 gophers && adduser -D -G gophers -u 1001 gopher
USER gopher
CMD ["/bin/dumbnailer"]
