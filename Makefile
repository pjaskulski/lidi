build-client:
	go build -o lidi-client ./cmd/client

build-server:
	go build -o lidi-server ./cmd/web

build-docker:
	docker build . -t lidi-server-alpine
