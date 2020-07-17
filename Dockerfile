FROM golang:alpine as builder

RUN apk add --no-cache make git
WORKDIR /nali-src
COPY . /nali-src
RUN go mod download && \
    make linux-amd64 && \
    mv ./bin/nali-linux-amd64 /nali

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /nali /
ENTRYPOINT ["/nali"]