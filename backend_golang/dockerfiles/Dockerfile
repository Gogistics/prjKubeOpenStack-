FROM golang:alpine

LABEL maintainer="Alan Tai <alan.tai@riverbed.com>"

ADD web.go ./
ADD apis /go/src/apis

# Update apk repositories
RUN echo "http://dl-3.alpinelinux.org/alpine/v3.8/main" > /etc/apk/repositories && \
  echo "http://dl-3.alpinelinux.org/alpine/v3.8/community" >> /etc/apk/repositories

# Install dependencies
RUN apk add --no-cache git dumb-init \
  && go get github.com/gorilla/websocket \
  && apk del git \
  && go build -o web web.go

# remove unnecessary stuff
RUN rm -rf /var/lib/apt/lists/* \
  /var/cache/apk/* \
  /usr/share/man \
  /tmp/*

ENV PORT 3000
EXPOSE 3000
ENTRYPOINT ["dumb-init", "./web"]