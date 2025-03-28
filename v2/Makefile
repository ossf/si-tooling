test:
	@echo "  >  Validating code ..."
	@go vet ./...
	@go clean -testcache
	@go test ./...

tidy:
	@echo "  >  Tidying go.mod ..."
	@go mod tidy

test-cov:
	@echo "Running tests and generating coverage output ..."
	@go test ./... -coverprofile coverage.out -covermode count
	@sleep 2 # Sleeping to allow for coverage.out file to get generated
	@echo "Current test coverage : $(shell go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+') %"