FROM golang:1.15 AS builder
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" ./cmd/api.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/completeworks.txt .
COPY --from=builder /app/static ./static
CMD ["/app/api"]
