FROM alpine:3.11

COPY ./webhook-server /run/server

EXPOSE 8443

ENTRYPOINT [ "/run/server" ]
