# Lag

## Get Started

```sh
docker-compose build
docker-compose up -d
```

```sh
curl http://localhost:8080/health
{"status": "ok"}
```

## API

Get a delayed response:
```sh
curl http://localhost:8080/?lag=10s
# 10s later
ok
```

Get 4xx and 5xx errors:
```sh
curl -i http://localhost:8080/error/418?lag=1s
HTTP/1.1 418 I'm a teapot

I'm a teapot
```

Download arbitrary amounts of data:
```sh
curl http://localhost:8080/data?s=125MB -o data
```

View all request headers
```sh
curl http://localhost:8080/headers
```
