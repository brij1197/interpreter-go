FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY app/ ./app/
RUN cd app && go build -o ../lox-interpreter *.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /build/lox-interpreter .
COPY web/ ./web/

EXPOSE 8080
CMD ["./lox-interpreter", "web"]