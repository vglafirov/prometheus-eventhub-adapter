FROM golang:latest

ENV BIN=prometheus-eventhub-adapter

COPY build/*-linux-amd64 /go/bin/$BIN

CMD /go/bin/$BIN