FROM ubuntu
RUN apt-get update && apt-get install rsyslog -y
COPY ./rsyslog.conf /etc/rsyslog.conf
ENTRYPOINT ["rsyslogd", "-n"]