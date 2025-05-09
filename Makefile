# Locally render the godocs site
pkgdocs:
	@echo "Launching pkgdocs ..."
	@cd v2 && \
		go run golang.org/x/pkgsite/cmd/pkgsite@latest -open
