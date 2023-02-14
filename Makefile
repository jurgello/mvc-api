run:
	docker run -p 8088:8088 myapp

build:
	docker build --tag myapp .

test:
	go test -v ./...