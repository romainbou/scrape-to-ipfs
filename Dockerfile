FROM golang:alpine AS builder

RUN apk add git

ADD ./ /go/src/github.com/romainbou/scrape-to-ipfs/

RUN set -ex && \
  cd /go/src/github.com/romainbou/scrape-to-ipfs && \       
  go get -v && \
  CGO_ENABLED=0 go build \
        -tags netgo \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  mv ./scrape-to-ipfs /usr/bin/scrape-to-ipfs

FROM busybox

COPY --from=builder /usr/bin/scrape-to-ipfs /usr/local/bin/scrape-to-ipfs

ENTRYPOINT [ "scrape-to-ipfs" ]