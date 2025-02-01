FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

# Сборка приложения
RUN apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/app .
COPY --from=builder /app/static/layouts /static/layouts

CMD ["./app"]