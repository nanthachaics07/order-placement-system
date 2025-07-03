run:
	go run cmd/main.go

tidy:
	go mod tidy

build:
	go build -o order-placement-system cmd/main.go

test:
	go test ./... -coverprofile=coverage.out

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean

build-up:
	docker-compose -f resources/docker/docker-compose.dev.yaml up -d

build-down:
	docker-compose -f resources/docker/docker-compose.dev.yaml down -v

gen-mock-order-processor:
	mockery \
	--name=OrderProcessorUseCase \
	--dir=internal/usecases/interfaces \
	--output=internal/mock/usecases \
	--outpkg=usecases
