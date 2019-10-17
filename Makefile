NAME := lag
TAG := $(shell git describe --long --match="v[0-9]*.[0-9]*.[0-9]*" --always)

.PHONY: build lint

build:
	docker build -t $(NAME):latest .
	docker tag $(NAME):latest $(NAME):$(TAG)

lint:
	go fmt ./...
