# Maintain semantic version
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
BUILD_NUMBER  ?= 0
PATCH_NUMBER  ?= 
VERSION_STRING = $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER)$(PATCH_NUMBER)

# Windows vs Unixlikes
ifdef COMSPEC
	DEL ?= del /s
else
	DEL ?= rm -r
endif

KUTTICMDFILES = cmd/kutti/*.go          \
				internal/pkg/cli/*.go   \
				internal/pkg/cmd/*.go   \
				internal/pkg/cmd/*/*.go \
				go.mod

# Targets
.PHONY: usage
usage:
	@echo "Usage: make windows|linux|clean"

out/kutti_linux_amd64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" cmd/kutti/*.go


out/kutti_windows_amd64.exe: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" cmd/kutti/*.go

.PHONY: windows
windows: out/kutti_windows_amd64.exe

.PHONY: linux
linux: out/kutti_linux_amd64

.PHONY: all
all: linux windows

.PHONY: clean
clean:
	$(DEL) out/