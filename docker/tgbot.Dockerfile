FROM golang:1.20

RUN apt-get update && \
    apt-get install -y gcc git pkg-config libssl-dev ca-certificates

COPY . /app

WORKDIR /app/cmd/tgbot

RUN go mod download

RUN go build -o tgbot

FROM ubuntu:20.04

RUN apt-get update && \
    apt-get install -y ca-certificates

COPY --from=0 /app/cmd/tgbot/tgbot /usr/local/bin/tgbot
COPY --from=0 /app/config/yaml/tgbot.yaml /home

EXPOSE 8081

CMD tgbot -config=/home/tgbot.yaml
CMD ["tgbot", "--config=/home/tgbot.yaml"]