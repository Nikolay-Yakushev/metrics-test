FROM golang:1.20-buster as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./application ./cmd/server/*.go


FROM alpine:3.17.2
WORKDIR /app
COPY --from=builder /app/application /app/application

CMD ["/app/application"]