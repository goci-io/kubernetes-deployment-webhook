FROM alpine:3.12

COPY bin/webhook-server /run/server

RUN addgroup runner && \
    mkdir -p /home/runner && \
    adduser -G runner runner --uid 1000 --disabled-password --home=/home/runner && \
    chown runner:runner /run/ /home/runner && \
    chmod g=u /run/ /home/runner

EXPOSE 9443
WORKDIR /run
USER runner

ENTRYPOINT [ "/run/server" ]
