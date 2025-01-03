#######################
## The builder image ##
#######################
FROM golang:1.23-alpine3.21 AS builder

ENV CGO_ENABLED=1

WORKDIR /app

RUN apk add --no-cache build-base=0.5-r3
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-w -s" -o rto cmd/rto/main.go

#######################
## The final image   ##
#######################
FROM alpine:3.21

LABEL maintainer="cmj <cmj@cmj.tw>"
LABEL description="The simple proxy service"
LABEL version="0.1.0"

WORKDIR /app

COPY --from=builder /app/rto /usr/local/bin/rto

CMD ["/usr/local/bin/rto", "-vv", "--cors"]
