.PHONY: build
build:
	go build .

.PHONY: docker
docker: build
	docker build -t dmon-service .

.PHONY: run
run: docker
	docker run -d dmon-service