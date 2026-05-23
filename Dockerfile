FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/devmate-api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S devmate && adduser -S devmate -G devmate

WORKDIR /app

COPY --from=builder /out/devmate-api /app/devmate-api

RUN mkdir -p /data && chown -R devmate:devmate /data /app

USER devmate

ENV PORT=8080
ENV DEVMATE_DB_PATH=/data/devmate.db

EXPOSE 8080

ENTRYPOINT ["/app/devmate-api"]
