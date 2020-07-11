FROM alpine:3.11

RUN apk update --no-cache && \
    apk upgrade --available

COPY packages.txt /etc/apk/packages.txt
RUN apk add --update --no-cache $(grep -v '^#' /etc/apk/packages.txt) && \
    rm -rf /var/cache/apk/* && \
    update-ca-certificates

RUN mkdir -p /run/workspace && \
    addgroup runner && \
    adduser --home /run/workspace --uid 1000 --disabled-password -G runner runner && \
    chown runner:runner /run/workspace/ && \
    chmod g=u /run/workspace/

WORKDIR /run/workspace
USER runner 

RUN mkdir ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    ssh-keyscan gitlab.com >> ~/.ssh/known_hosts

ENTRYPOINT [ "/bin/bash", "-c" ]
