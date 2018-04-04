build:
	go build -o hermes ./cmd/server.go
	go build ./cmd/publisher.go
	go build ./cmd/subscriber.go
