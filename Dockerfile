FROM scratch

ADD main ./

EXPOSE 2448:2448/udp
CMD ["/main"]
