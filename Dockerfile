FROM golang:1.15.4-alpine as build

WORKDIR /go/src/github.com/Bhinneka/user-service/
ADD . /go/src/github.com/Bhinneka/user-service
COPY go.mod go.sum .env /go/src/github.com/Bhinneka/user-service/

ARG SSH_PRIVATE_KEY
ENV BUILD_PACKAGES="git curl build-base make openssh"
ENV GO111MODULE=on

RUN apk update && apk add --no-cache $BUILD_PACKAGES \
      && mkdir /root/.ssh/ && mv id_rsa /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa && touch /root/.ssh/known_hosts \
      && mkdir -p /usr/filebeat/ && cp -r _filebeat/* /usr/filebeat/ \
      && ssh-keyscan github.com >> /root/.ssh/known_hosts \
      && git config --global url."git@github.com:".insteadOf "https://github.com/" \
      && make user-service-linux \
      && apk del $BUILD_PACKAGES

FROM alpine:3.13.5
RUN apk update \
      && apk add rsyslog \
      && apk add supervisor \
      && apk add tzdata
ARG BINARY_PATH=/go/src/github.com/Bhinneka/user-service
RUN mkdir -p /usr/filebeat/ && mkdir -p /var/log/

ADD _build/rsyslog.conf /etc/rsyslog.conf
ADD _build/rsyslog.d/ /etc/rsyslog.d/
ADD _build/supervisord.conf /etc/supervisord.conf

RUN chmod 644 /etc/rsyslog.d/ && chmod 644 /etc/rsyslog.conf && cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

VOLUME ["/var/log"]
COPY --from=build /usr/filebeat/* /usr/filebeat/
COPY --from=build $BINARY_PATH/.env $BINARY_PATH/.env
COPY --from=build $BINARY_PATH/config/rsa/* $BINARY_PATH/config/rsa/
COPY --from=build $BINARY_PATH/schema/json/* $BINARY_PATH/schema/json/

USER root
RUN chmod go-w /usr/filebeat/filebeat.yml
VOLUME [ "/usr/filebeat/" ]

EXPOSE 8082
EXPOSE 8081

COPY --from=build $BINARY_PATH/user-service-linux $BINARY_PATH/user-service-linux

ENTRYPOINT ["sh", "-c", "supervisord -nc /etc/supervisord.conf"]
