FROM golang:1.20

RUN apt-get update && \
    apt-get install -y gcc git pkg-config libssl-dev ca-certificates

COPY . /app

WORKDIR /app/cmd/passman

RUN go mod download

RUN go build -o passman

FROM ubuntu:20.04

RUN apt-get update && \
    apt-get install -y ca-certificates

COPY --from=0 /app/cmd/passman/passman /usr/local/bin/passman
COPY --from=0 /app/config/yaml/passman.yaml /home

EXPOSE 8080

CMD ["passman", "--config=/home/passman.yaml"]