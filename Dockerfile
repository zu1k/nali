FROM golang:1.27-alpine AS builder

RUN apk add --no-cache make git
WORKDIR /nali-src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make docker && \
    mv ./bin/nali-docker /nali

FROM alpine:3

RUN apk add --no-cache ca-certificates
COPY --from=builder /nali /
ENTRYPOINT ["/nali"]
