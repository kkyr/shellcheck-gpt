.PHONY: release
release:
	./build/release.sh

.PHONY: lint
lint:
	./build/lint.sh ./...