# Maintain semantic version
# Also change in cmd/kutti/main.go
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 5
BUILD_NUMBER  ?= 0
PATCH_NUMBER  ?= 
VERSION_STRING = $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER)$(PATCH_NUMBER)

KUTTICMDFILES = cmd/kutti/*.go          \
				internal/pkg/cli/*.go   \
				internal/pkg/cmd/*.go   \
				internal/pkg/cmd/*/*.go \
				go.mod \
				Makefile

# Targets
.PHONY: usage
usage:
	@echo "Usage: make linux|windows|mac|mac-intel|linux-install-script|windows-installer|mac-install-script|mac-intel-install-script|all|installers|clean"

out/:
	mkdir out

out/kutti_linux_amd64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/get-kutti-linux-amd64.sh: build/package/posix-install-script/generate-script.sh out/
	CURRENT_VERSION=${VERSION_STRING} GOOS=linux GOARCH=amd64 $< > $@

cmd/kutti/rsrc_windows_amd64.syso: cmd/kutti/winres/*
	go-winres make --in=cmd/kutti/winres/winres.json --out=cmd/kutti/rsrc --arch=amd64 --product-version=${VERSION_STRING} --file-version=${VERSION_STRING}

out/kutti_windows_amd64.exe: $(KUTTICMDFILES) cmd/kutti/rsrc_windows_amd64.syso
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/kutti-windows-installer.exe: build/package/kutti-windows-installer/kutti-windows-installer.nsi out/kutti_windows_amd64.exe
	makensis -NOCD -V3 -- $<

out/kutti_darwin_amd64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/get-kutti-darwin-amd64.sh: build/package/posix-install-script/generate-script.sh out/
	CURRENT_VERSION=${VERSION_STRING} GOOS=darwin GOARCH=amd64 $< > $@

out/kutti_darwin_arm64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/get-kutti-darwin-arm64.sh: build/package/posix-install-script/generate-script.sh out/
	CURRENT_VERSION=${VERSION_STRING} GOOS=darwin GOARCH=arm64 $< > $@

.PHONY: linux
linux: out/kutti_linux_amd64

.PHONY: linux-install-script
linux-install-script: out/get-kutti-linux-amd64.sh

.PHONY: windows
windows: out/kutti_windows_amd64.exe

.PHONY: windows-installer
windows-installer: out/kutti-windows-installer.exe

.PHONY: mac-intel
mac-intel: out/kutti_darwin_amd64

.PHONY: mac-intel-install-script
mac-intel-install-script: out/get-kutti-darwin-amd64.sh

.PHONY: mac
mac: out/kutti_darwin_arm64

.PHONY: mac-install-script
mac-install-script: out/get-kutti-darwin-arm64.sh

.PHONY: all
all: linux windows mac mac-intel

.PHONY: installers
installers: linux-install-script windows-installer mac-install-script mac-intel-install-script

.PHONY: resourceclean
resourceclean:
	rm -f cmd/kutti/rsrc_windows_amd64.syso

.PHONY: binclean
binclean:
	rm -r -f out/

.PHONY: clean
clean: resourceclean binclean 

out/man:
	mkdir -p out/man

.PHONY: manpagedocs
manpagedocs: out/man $(KUTTICMDFILES)
	go run internal/cmd/gendoc/main.go -o out/man -t manpages

out/markdown:
	mkdir -p out/markdown

.PHONY: markdowndocs
markdowndocs: out/markdown $(KUTTICMDFILES)
	go run internal/cmd/gendoc/main.go -o out/markdown -t markdown

.PHONY: manpagedocsclean
manfiledocsclean:
	rm -rf out/man

.PHONY: markdowndocsclean
markdowndocsclean:
	rm -rf out/markdown

.PHONY: docsclean
docsclean: manfiledocsclean markdowndocsclean
