FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# Default environment variables
ENV PORT=3000
EXPOSE 3000

CMD ["./main"]
