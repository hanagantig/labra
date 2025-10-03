FROM golang:1.25.1-alpine3.22 AS build

ENV APP_HOME /opt/build
WORKDIR $APP_HOME

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o labra

FROM alpine:3.22

WORKDIR /opt/bin

#RUN apk update \
#    && apk add ca-certificates tzdata \
#    && ln -fns /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
#    && rm -rf /tmp/* /var/tmp/* \

COPY --from=build /opt/build/labra .

CMD ["/opt/bin/labra", "http", "--config=./config.yaml"]
