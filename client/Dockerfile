
FROM golang:alpine AS builder
WORKDIR /app
COPY go.* .
COPY communication ./communication
COPY utils ./utils
COPY client ./client
COPY api ./api

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./client

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY database ./database
COPY logs ./logs
CMD ["./main"]
