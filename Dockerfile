# Build stage
FROM golang:1.14-buster as builder

ENV GO111MOD=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Release stage
FROM alpine as release

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates nginx \
        && update-ca-certificates 2>/dev/null || true

ENV APP_ENV=release

# Nginx config
COPY config/nginx/app.conf.template /etc/nginx/conf.d/default.conf

COPY --from=builder /app/backend /app/

EXPOSE $PORT

CMD ./bin/start.sh
