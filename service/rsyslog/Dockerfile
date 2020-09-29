FROM ubuntu


RUN apt-get update && apt install rsyslog -y

# RUN apk update \
#     && apk add rsyslog --no-cache \
#     && rm -rf /var/cache/apk/*
COPY ./rsyslog.conf /etc/rsyslog.conf
ENTRYPOINT ["rsyslogd", "-n"]
