FROM golang:1.16-alpine as build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... && \
    go install -v ./...

FROM alpine as runtime
COPY --from=build /go/bin/cron-viewer /usr/local/bin/
ENTRYPOINT ["cron-viewer"]