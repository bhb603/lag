# Lag

A server that lags (on purpose).

## Get Started

```sh
docker-compose build
docker-compose up -d
```

```sh
curl http://localhost:8080/health
OK
```

## API

Make an HTTP request with a defined response time:
```sh
curl http://localhost:8080/?t=10s
# 10s later
OK
```

Make an HTTP request which responds with 4xx or 5xx errors:
```sh
curl -i http://localhost:8080/error/418?t=1s
HTTP/1.1 418 I'm a teapot

I'm a teapot
```
