FROM golang:1.21.5-alpine3.17 as builder
WORKDIR /src
COPY go.* ./
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go mod download
RUN go build -o walletAPI ./cmd/main.go

FROM alpine:3.17 as binary
COPY --from=builder /src/walletAPI .
EXPOSE 3000
CMD ["/walletAPI"]