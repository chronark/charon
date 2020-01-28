FROM alpine:3.11

RUN  apk add --update rsyslog \
  && rm -rf /var/cache/apk/*

EXPOSE 514 514/udp

ENTRYPOINT [ "rsyslogd", "-n" ]

# VOLUME [ "/var/log", "/etc/rsyslog.d" ]

# for some reason, the apk comes built with a v5
# config file. using this one for v8:
# COPY ./etc/rsyslog.conf /etc/rsyslog.conf