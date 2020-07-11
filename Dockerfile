FROM scratch

COPY bin/webhook-server /run/server

EXPOSE 8443

ENTRYPOINT [ "/run/server" ]
