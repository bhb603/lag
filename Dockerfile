FROM golang:1.15.2 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/lag main.go

FROM alpine:latest
COPY --from=builder /bin/lag /bin/lag
COPY config.yaml /etc/lag/config.yaml
CMD ["/bin/lag", "serve"]
