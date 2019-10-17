NAME := bhb603/lag
TAG := $(shell git describe --match="[0-9]*.[0-9]*.[0-9]*")

.PHONY: build push lint

build:
	docker build -t $(NAME):$(TAG) .
	docker tag $(NAME):$(TAG) $(NAME):latest

push:
	docker push $(NAME):$(TAG)
	docker push $(NAME):latest

lint:
	go fmt ./...
