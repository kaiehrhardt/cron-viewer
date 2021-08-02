FROM golang:1.16-alpine as build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... && \
    go install -v ./... 

FROM alpine as runtime
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime && \
    echo "Europe/Berlin" > /etc/timezone
COPY --from=build /go/bin/cron-viewer /usr/local/bin/
ENTRYPOINT ["cron-viewer"]