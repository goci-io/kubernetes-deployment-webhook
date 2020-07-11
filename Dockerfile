FROM scratch

COPY server_linux_amd64/bin/webhook-server /run/server

EXPOSE 8443

ENTRYPOINT [ "/run/server" ]
