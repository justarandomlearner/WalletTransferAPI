FROM golang:1.21.5-alpine3.17 as base
WORKDIR /src/walletAPI
ADD . . 
RUN go mod download
RUN go build -o walletAPI ./cmd

FROM alpine:3.17 as binary
WORKDIR /src/app
COPY --from=base /src/walletAPI/walletAPI .
EXPOSE 3000
CMD ["/src/app/walletAPI"]