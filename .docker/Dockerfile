FROM golang:1.25.1-alpine3.22 AS build

ENV APP_HOME /opt/build
COPY . $APP_HOME

WORKDIR $APP_HOME

RUN go build

FROM alpine:3.22

ENV BUILD_DIR /opt/build
ENV APP_HOME /opt/bin

WORKDIR $APP_HOME

#RUN apk update \
#    && apk add ca-certificates tzdata \
#    && ln -fns /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
#    && rm -rf /tmp/* /var/tmp/* \

COPY --from=build $BUILD_DIR/labra .

CMD ["/opt/bin/labra", "forms"]
