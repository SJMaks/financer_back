FROM golang:1.26.0-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/financer/exe /app/cmd/financer

FROM alpine:3.24.2
WORKDIR /app
COPY --from=builder /app/cmd/financer/exe /app/
CMD ["/app/exe"]