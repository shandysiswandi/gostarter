FROM alpine:3.21.0
ENV GOOSE_VERSION_TAG="v3.23.0"
RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/*
RUN wget -O /bin/goose \
    https://github.com/pressly/goose/releases/download/$GOOSE_VERSION_TAG/goose_linux_x86_64 
RUN chmod +x /bin/goose

WORKDIR /migrations