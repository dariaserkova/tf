FROM golang:1.12.5-stretch as builder

COPY *.go ./src/app/

WORKDIR /go/src/app

RUN go get -v -d -t ./... && sync && \
    CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-w' -installsuffix cgo -o boltsrv .
# https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671

FROM alpine:3.9

RUN apk --no-cache add ca-certificates

RUN addgroup --gid 1024 app && adduser -h /app --disabled-password -u 1024 -G app app
USER app
WORKDIR /app
EXPOSE 5000
COPY --from=builder /go/src/app/boltsrv .
CMD [ "./boltsrv" ]