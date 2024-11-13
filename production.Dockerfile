FROM golang:1.19

WORKDIR /go/src/app
COPY . .
COPY config-docker.json config.json

RUN go mod download
RUN env GOOS=linux GOARCH=amd64 go build -o miniwebserver

CMD ["/go/src/app/miniwebserver","config.json"]