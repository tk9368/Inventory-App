FROM golang:1.21-alpine AS base
WORKDIR /src
RUN go mod init inventory-app &&go get github.com/lib/pq
COPY main.go .
RUN go build -o store-app .
FROM alpine:latest
WORKDIR /app
COPY --from=base /src/store-app .
EXPOSE 8080
CMD ["./store-app"]
