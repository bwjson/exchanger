FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /exchanger ./cmd/app

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /exchanger .
COPY config config
USER nobody
CMD ["./exchanger"]
