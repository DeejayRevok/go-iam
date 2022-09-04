FROM golang:1.18.3

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
