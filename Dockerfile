# Build stage
FROM golang:1.23 as builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git \
    bash \
    build-essential \
    libc-dev \
    libatk1.0-0 libatk-bridge2.0-0 libcups2 libxcomposite1 \
    libxdamage1 libxrandr2 libgbm1 libasound2 libpangocairo-1.0-0 \
    libgtk-3-0 libnss3 libxshmfence1 \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go mod tidy
RUN go run main.go -download-only
RUN go build -o app .

FROM debian:stable-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    libatk1.0-0 libatk-bridge2.0-0 libcups2 libxcomposite1 \
    libxdamage1 libxrandr2 libgbm1 libasound2 libpangocairo-1.0-0 \
    libgtk-3-0 libnss3 libxshmfence1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /root/.cache/rod /root/.cache/rod
COPY --from=builder /app/app .
COPY --from=builder /app/views /app/views

EXPOSE 3000

CMD ["./app"]
