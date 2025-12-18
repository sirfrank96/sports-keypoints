FROM golang:1.25.3 AS builder

WORKDIR /app/server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go-server .

EXPOSE 50052

CMD ["./go-server"]
