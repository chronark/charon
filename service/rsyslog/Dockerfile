FROM ubuntu:18.04
RUN apt-get update && apt-get install rsyslog=8.32.0 -y
COPY ./rsyslog.conf /etc/rsyslog.conf
ENTRYPOINT ["rsyslogd", "-n"]
