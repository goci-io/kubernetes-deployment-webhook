FROM alpine:3.11

RUN apk update --no-cache && \
    apk upgrade --available

COPY packages.txt /etc/apk/packages.txt
RUN apk add --update --no-cache $(grep -v '^#' /etc/apk/packages.txt) && \
    rm -rf /var/cache/apk/* && \
    update-ca-certificates

RUN mkdir -p /run/workspace && \
    addgroup runner && \
    adduser -S -G runner runner && \
    chown runner:runner /run/workspace/ && \
    chmod g=u /run/workspace/

WORKDIR /run/workspace

ENTRYPOINT [ "/bin/bash", "-c" ]
