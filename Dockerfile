FROM golang:1.16-alpine
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... && \
    go install -v ./...
ENTRYPOINT ["cron-viewer"]