FROM golang:1.22.4-bullseye

WORKDIR /go/src/app

COPY go.mod .
COPY main.go .

RUN go get

RUN go build

CMD /go/src/app/prom-counter
