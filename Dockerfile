# 多階段構建 - 後端
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# 安裝構建依賴
RUN apk add --no-cache git

# 複製 go mod 文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 複製後端源碼
COPY backend/ ./

# 構建後端
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tspl-simulator .

# 多階段構建 - 前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 複製 package 文件
COPY frontend/package*.json ./
RUN npm ci

# 複製前端源碼
COPY frontend/ ./

# 構建前端
RUN npm run build

# 最終運行階段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從構建階段複製後端二進制文件
COPY --from=backend-builder /app/backend/tspl-simulator .

# 從構建階段複製前端構建文件
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# 複製環境變數範例
COPY backend/.env.example ./.env.example

# 創建 data 目錄
RUN mkdir -p data/API_print data/MQTT_print

# 暴露端口
EXPOSE 8080

# 運行後端
CMD ["./tspl-simulator"]
