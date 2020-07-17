.PHONY: build
build:
	go build .

.PHONY: docker
docker: build
	docker build -t dmon-service .

.PHONY: run
run: docker
	docker run --name dmon-service -d dmon-service

.PHONY: stop
stop:
	docker rm --force dmon-service