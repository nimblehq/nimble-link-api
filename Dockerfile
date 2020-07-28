# Build stage
FROM golang:1.14-buster as builder

ENV GO111MOD=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Release stage
FROM scratch as release

ENV APP_ENV=release

COPY --from=builder /app/backend /app/

EXPOSE 8080

ENTRYPOINT ["/app/backend"]
