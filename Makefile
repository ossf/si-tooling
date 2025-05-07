# Locally render the godocs site
pkgdocs:
	@echo "Launching pkgdocs ..."
	@cd v2 && \
		go run golang.org/x/pkgsite/cmd/pkgsite@latest -open

test-cov:
	@echo "Running tests and generating coverage output ..."
	@cd v2 && \
		go test ./... -coverprofile coverage.out -covermode count
	@sleep 2 # Sleeping to allow for coverage.out file to get generated

covcheck: test-cov
	@COVERAGE=$(shell cd v2 && go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'); \
	THRESHOLD=75.0; \
	echo "Test coverage: $$COVERAGE%"; \
	echo "Coverage threshold: $$THRESHOLD%"; \
	if [ $$(echo "$$COVERAGE < $$THRESHOLD" | bc) -gt 0 ]; then \
		echo "WARNING: Test coverage ($$COVERAGE%) is below the threshold ($$THRESHOLD%)!"; \
		exit 1; \
	else \
		echo "Test coverage ($$COVERAGE%) exceeds the threshold ($$THRESHOLD%)."; \
		exit 0; \
	fi

PHONY: test-cov covcheck pkgdocs
