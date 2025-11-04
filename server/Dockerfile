FROM golang:alpine AS builder
WORKDIR /app
COPY go.* .
COPY server ./server
COPY communication ./communication
COPY models ./models
COPY utils ./utils
COPY api ./api

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./server

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY database ./database
COPY logs ./logs
CMD ["./main"]
