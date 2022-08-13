FROM golang:1.18.4-alpine3.15

WORKDIR /go/src/websocket-in-go-boilerplate

COPY . /go/src/websocket-in-go-boilerplate

RUN go mod download

RUN go install

RUN go build .

EXPOSE 8000

CMD ["./websocket-in-go-boilerplate"] --v
