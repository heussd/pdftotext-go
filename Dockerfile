FROM golang:1.22.0-alpine3.18

RUN   apk add --no-cache --update poppler-utils

WORKDIR /app

COPY testdata ./testdata
COPY go.* ./
COPY *.go ./

RUN go test -v