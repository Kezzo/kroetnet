# first run: CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
FROM scratch

ADD main ./

EXPOSE 2448:2448/udp
CMD ["/main"]
