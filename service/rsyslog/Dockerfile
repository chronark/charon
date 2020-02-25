FROM ubuntu:18.04
RUN apt-get update \
    && apt-get install rsyslog=8.32.0 -y --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
    
COPY ./rsyslog.conf /etc/rsyslog.conf
ENTRYPOINT ["rsyslogd", "-n"]
