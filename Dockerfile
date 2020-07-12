FROM alpine:3.12

COPY bin/webhook-server /run/server

RUN addgroup runner && \
    adduser -G runner runner --uid 1000 --disabled-password --home=/run && \
    chown runner:runner /run/ && \
    chmod g=u /run/

EXPOSE 8443
WORKDIR /run
USER runner

ENTRYPOINT [ "/run/server" ]
