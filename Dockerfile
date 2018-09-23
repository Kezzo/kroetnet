FROM golang:alpine

WORKDIR /go/src/kroetnet
COPY . .

RUN go install -v ./...

EXPOSE 2448:2448/udp
CMD ["kroetnet"]
