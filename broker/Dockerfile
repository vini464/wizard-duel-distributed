FROM golang:alpine AS builder
WORKDIR /app
COPY go.* .
COPY communication ./communication
COPY broker ./broker
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main broker/main.go broker/broker.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
