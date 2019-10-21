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

Get a delayed response:
```sh
curl http://localhost:8080/?t=10s
# 10s later
OK
```

Get 4xx and 5xx errors:
```sh
curl -i http://localhost:8080/error/418?t=1s
HTTP/1.1 418 I'm a teapot

I'm a teapot
```

Download arbitrary amounts of data:
```sh
curl http://localhost:8080/data?s=125MB -o data
```
