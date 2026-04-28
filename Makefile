# Service-specific configuration
SERVICE_NAME := profile
APP_DIRS     := apps/default apps/devices apps/geolocation apps/settings

# Bootstrap: download shared Makefile.common if missing
ifeq (,$(wildcard .tmp/Makefile.common))
  $(shell mkdir -p .tmp && curl -sSfL https://raw.githubusercontent.com/antinvestor/common/main/Makefile.common -o .tmp/Makefile.common)
endif

include .tmp/Makefile.common

# Dart proto modules — each gets its own buf.gen.dart.<module>.yaml so that
# generation is scoped to a single module and dart packages don't leak each
# other's types. See proto/buf.gen.dart.*.yaml.
DART_MODULES := profile device settings ocr geolocation

.PHONY: proto-generate-dart
proto-generate-dart: $(BIN)/buf ## Regenerate the per-module dart SDKs
	@if [ ! -d "$(PROTO_DIR)" ]; then exit 0; fi
	@# Purge sibling-module dirs first (these are buf-generate noise that
	@# leaks when plugins are run unscoped against the workspace). Hand-
	@# written files like client.dart sit in the same lib/src/ root so we
	@# explicitly target only the sibling-module subdirs.
	@for pkg in $(DART_MODULES); do \
		for other in $(DART_MODULES); do \
			[ "$$pkg" = "$$other" ] && continue; \
			rm -rf sdk/dart/$$pkg/lib/src/$$other; \
		done; \
	done
	@for m in $(DART_MODULES); do \
		echo "==> dart $$m"; \
		(cd $(PROTO_DIR) && buf generate --template buf.gen.dart.$$m.yaml $$m); \
	done

# Wire dart generation into the standard proto-generate pipeline so that
# `make proto-generate` produces both Go/OpenAPI and dart SDKs in one step.
proto-generate: proto-generate-dart

format: ## Format Go files (used by pre-commit hook)
	gofmt -w .
