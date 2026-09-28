FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o calendar-server .

FROM alpine:3.24

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/calendar-server .

EXPOSE 8080

CMD ["./calendar-server"]
