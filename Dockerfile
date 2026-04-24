# ==========================================
# STAGE 1: BUILDER
# ==========================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_TARGET
ARG BINARY_NAME=app_service

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/${BINARY_NAME} ${BUILD_TARGET}

# ==========================================
# STAGE 2: PRODUCTION RUNNER
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=builder /app/app_service .

COPY etc/ ./etc/

COPY --from=builder /app/*.pb .

CMD ["./app_service"]