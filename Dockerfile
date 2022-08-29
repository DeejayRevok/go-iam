FROM golang:1.18.3

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go install

ARG DEVELOPMENT
ENV DEVELOPMENT ${DEVELOPMENT}
CMD ["sh", "run.sh"]
