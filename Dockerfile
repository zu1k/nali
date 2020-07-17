FROM golang:alpine as builder

RUN apk add --no-cache make git
WORKDIR /nali-src
COPY . /nali-src
RUN go mod download && \
    make docker && \
    mv ./bin/nali-docker /nali

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /nali /
ENTRYPOINT ["/nali"]