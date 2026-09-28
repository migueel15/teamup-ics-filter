FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o calendar-server .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/calendar-server .

EXPOSE 8080

CMD ["./calendar-server"]
