# Build Stage
FROM golang:1.23 AS builder

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
RUN go build -o app .

# Runtime Stage
FROM frolvlad/alpine-glibc:alpine-3.17

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    dumb-init \
    bash \
    ca-certificates \
    && mkdir -p /root/.cache

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/views /app/views

EXPOSE 3000
ENTRYPOINT ["dumb-init", "./app"]
