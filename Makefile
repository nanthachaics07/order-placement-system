run:
	go run cmd/main.go

tidy:
	go mod tidy

build:
	go build -o order-placement-system cmd/main.go

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -n 1


test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | tee /dev/tty | tail -n 1
	@total=$$(go tool cover -func=coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); \
	threshold=80.0; \
	echo "Total Coverage: $$total%"; \
	if [ $$(echo "$$total < $$threshold" | bc) -eq 1 ]; then \
		echo "Coverage MISS < $$threshold%"; exit 1; \
	else \
		echo "Coverage PASS >= $$threshold%"; fi

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
