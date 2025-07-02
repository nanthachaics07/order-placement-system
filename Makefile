run:
	go run cmd/main.go

build:
	go build -o order-placement-system cmd/main.go

test:
	go test ./... -v

clean:
	go clean
