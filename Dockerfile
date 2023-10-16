FROM golang:1.21 AS builder

WORKDIR /usr/src/swissbank

COPY src/ .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -v -o /bin/swissbank cmd/http/main.go

FROM alpine

COPY --from=builder /bin/swissbank /bin/swissbank

RUN apk --no-cache add tzdata ca-certificates

ENTRYPOINT /bin/swissbank
