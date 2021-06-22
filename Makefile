# Maintain semantic version
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
BUILD_NUMBER  ?= 0
PATCH_NUMBER  ?= 
VERSION_STRING = $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER)$(PATCH_NUMBER)

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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/kutti_windows_amd64.exe: $(KUTTICMDFILES) cmd/kutti/winres/*
	go-winres make --in=cmd/kutti/winres/winres.json --out=cmd/kutti/rsrc --arch=amd64 --product-version=${VERSION_STRING} --file-version=${VERSION_STRING}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

.PHONY: windows
windows: out/kutti_windows_amd64.exe

.PHONY: linux
linux: out/kutti_linux_amd64

.PHONY: all
all: linux windows

.PHONY: resourceclean
resourceclean: cmd/kutti/*.syso
	rm cmd/kutti/*.syso

.PHONY: binclean
binclean: out/*
	rm -r out/

.PHONY: clean
clean: binclean resourceclean
