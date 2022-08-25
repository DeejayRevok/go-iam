FROM golang:1.18.3

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go install

EXPOSE 8888
CMD ["go", "run", "main.go"]