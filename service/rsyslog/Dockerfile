FROM alpine:3.12
RUN apk update \
    && apk add rsyslog=8.1911.0-r1 --no-cache \
    && rm -rf /var/cache/apk/*
COPY ./rsyslog.conf /etc/rsyslog.conf
ENTRYPOINT ["rsyslogd", "-n"]
