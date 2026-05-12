# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
#RUN go build -o main .
RUN go build -o main ./cmd/main.go

# Stage 2: Runtime (image jauh lebih kecil ~10MB vs ~800MB)
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

# Dockerfile ini untuk production

# docker build -t go-fiber-app:latest .
# # save image jadi file
# docker save go-fiber-app:latest | gzip > go-fiber-app_latest.tar.gz

# copy file ke server
# scp go-fiber-app_latest.tar.gz user@server:/path/to/deploy

# # di server, load image dari file
# cd /path/to/deploy
# docker load < go-fiber-app_latest.tar.gz
# docker images