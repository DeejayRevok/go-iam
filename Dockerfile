FROM golang:1.18.3

RUN wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
RUN apt-get install apt-transport-https
RUN echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
RUN apt-get update && apt-get install filebeat

COPY ./tools/filebeat.yml /etc/filebeat/filebeat.yml

RUN mkdir /var/log/uaa
RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go install

ARG DEVELOPMENT
ENV DEVELOPMENT ${DEVELOPMENT}
CMD ["sh", "run.sh"]
