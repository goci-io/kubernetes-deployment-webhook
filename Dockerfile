FROM scratch

COPY ./webhook-server /run/server

EXPOSE 8443

ENTRYPOINT [ "/run/server" ]
