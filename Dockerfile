FROM golang:1.18.4-alpine3.15

WORKDIR /go/src/websocket-in-go-example

COPY . /go/src/websocket-in-go-example

RUN go mod download

RUN go install

RUN go build .

EXPOSE 8000

CMD ["./websocket-in-go-example"] --v
