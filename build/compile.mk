#
# Makefile Fragment for Compiling
#

GO         ?= go
BUILD_DIR  ?= ./bin
PROJECT_MODULE ?= $(shell $(GO) list -m)
# $b replaced by the binary name in the compile loop, -s/w remove debug symbols
LDFLAGS    ?= "-s -w -X main.appName=$$b"
SRCDIR     ?= .
COMPILE_OS ?= darwin linux windows

# Determine commands by looking into cmd/*
COMMANDS   ?= $(wildcard ${SRCDIR}/cmd/*)

# Determine binary names by stripping out the dir names
BINS       := $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

compile-clean:
	@echo "=== $(PROJECT_NAME) === [ compile-clean    ]: removing binaries..."
	@rm -rfv $(BUILD_DIR)/*

compile: deps compile-only

compile-all: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	@for b in $(BINS); do \
		for os in $(COMPILE_OS); do \
			BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
			if [ "$$os" = "windows" ]; then \
				echo "=== $(PROJECT_NAME) === [ compile          ][$(GOARCH)]:     $(BUILD_DIR)/$$os/$$b.exe"; \
				CGO_ENABLED=0 GOOS=$$os $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$$os/$$b.exe $$BUILD_FILES ; \
			else \
				echo "=== $(PROJECT_NAME) === [ compile          ][$(GOARCH)]:     $(BUILD_DIR)/$$os/$$b"; \
				CGO_ENABLED=0 GOOS=$$os $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$$os/$$b $$BUILD_FILES ; \
			fi \
		done \
	done

compile-only: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ compile          ][$(GOARCH)]:      $(BUILD_DIR)/$(GOOS)/$$b"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(GOOS)/$$b $$BUILD_FILES ; \
	done

# Override GOOS for these specific targets
compile-darwin: GOOS=darwin
compile-darwin: deps-only compile-only

compile-linux: GOOS=linux
compile-linux: deps-only compile-only

compile-linux-x86: GOOS=linux
compile-linux-x86: GOARCH=386
compile-linux-x86: deps-only compile-only

compile-windows: GOOS=windows
compile-windows: deps-only compile-only

compile-windows-x86: GOOS=windows
compile-windows-x86: GOARCH=386
compile-windows-x86: deps-only compile-only

.PHONY: clean-compile compile compile-darwin compile-linux compile-only compile-windows
