run:
	go run cmd/main.go

build:
	go build -o order-placement-system cmd/main.go

test:
	go test ./... -coverprofile=coverage.out

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean
