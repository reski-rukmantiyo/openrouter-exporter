FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/openrouter-exporter .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/openrouter-exporter /usr/local/bin/openrouter-exporter

EXPOSE 9837
ENTRYPOINT ["openrouter-exporter"]
