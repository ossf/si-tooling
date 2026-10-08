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

PHONY: test-cov covcheck pkgdocs check-spec-examples

# Confirm the spec examples in test_data match the published Security Insights release
SPEC_VERSION ?= v2.2.0
check-spec-examples:
	@echo "Comparing test_data/spec-$(SPEC_VERSION) with the Security Insights $(SPEC_VERSION) examples ..."
	@for f in v2/si/test_data/spec-$(SPEC_VERSION)/*.yml; do \
		curl -sfL "https://raw.githubusercontent.com/ossf/security-insights/$(SPEC_VERSION)/examples/$$(basename $$f)" | diff -q - "$$f" >/dev/null \
			|| { echo "$$f differs from the $(SPEC_VERSION) example"; exit 1; }; \
	done
	@echo "All spec examples match."
