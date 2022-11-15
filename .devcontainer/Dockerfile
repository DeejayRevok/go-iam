FROM golang:1.18.3

RUN wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
RUN apt-get install apt-transport-https
RUN echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
<<<<<<< Updated upstream
RUN apt-get update && apt-get install filebeat

RUN wget -qO- https://repos.influxdata.com/influxdb.key | apt-key add -
RUN apt-get install -y lsb-release
RUN echo "deb https://repos.influxdata.com/debian $(lsb_release -cs) stable"| tee /etc/apt/sources.list.d/influxdb.list
RUN apt-get update && apt-get install telegraf

COPY ./tools/filebeat.yml /etc/filebeat/filebeat.yml
COPY ./tools/telegraf.conf /etc/telegraf/telegraf.conf
=======
RUN apt-get update && apt-get install filebeat && apt-get install metricbeat=7.11.2
COPY ./tools/filebeat.yml /etc/filebeat/filebeat.yml
COPY ./tools/metricbeat.yml /etc/metricbeat/metricbeat.yml
>>>>>>> Stashed changes

RUN mkdir /var/log/uaa
RUN mkdir /app
COPY . /app
WORKDIR /app

ENV CGO_ENABLED=1

RUN go install

RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go install -v github.com/vektra/mockery/v2@latest

ARG DEVELOPMENT
ENV DEVELOPMENT ${DEVELOPMENT}
CMD ["sh", "run.sh"]
