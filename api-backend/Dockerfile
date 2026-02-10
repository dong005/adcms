# 构建阶段
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o adcms-server ./cmd/main.go

# 运行阶段
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/adcms-server .
COPY --from=builder /app/config.yaml .

RUN mkdir -p uploads logs

EXPOSE 8004

CMD ["./adcms-server"]
