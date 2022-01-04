FROM golang:alpine AS builder

RUN apk update

RUN apk add --no-cache \
        libc6-compat git

WORKDIR /garlicshare

RUN go get github.com/R4yGM/garlicshare


FROM alpine:edge
COPY --from=builder /go/bin/garlicshare /bin/garlicshare
RUN apk add --no-cache tor

ENTRYPOINT ["garlicshare"]