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
COPY --from=build /go/src/app/config.yml /etc/cron-viewer/config.yml
EXPOSE 8080
ENTRYPOINT ["cron-viewer", "-config", "/etc/cron-viewer/config.yml"]