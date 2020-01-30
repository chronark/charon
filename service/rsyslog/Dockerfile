FROM ubuntu:18.04
RUN apt-get update && apt-get install rsyslog -y --no-install-recommends
COPY ./rsyslog.conf /etc/rsyslog.conf
ENTRYPOINT ["rsyslogd", "-n"]
