# Maintain semantic version
# Also change in cmd/kutti/main.go
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 6
BUILD_NUMBER  ?= 0
PATCH_NUMBER  ?= 
VERSION_STRING = $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER)$(PATCH_NUMBER)
VERSION_FOUR_PART = $(VERSION_MAJOR).$(VERSION_MINOR).$(BUILD_NUMBER).0

KUTTICMDFILES = cmd/kutti/*.go          \
				internal/pkg/cli/*.go   \
				internal/pkg/cmd/*.go   \
				internal/pkg/cmd/*/*.go \
				go.mod \
				Makefile

# Targets
.PHONY: usage
usage:
	@echo "Usage: make linux|windows|mac|mac-intel|deb|rpm|archlinux|linux-packages|completions|windows-installer|all|installers|clean"

out/:
	mkdir -p out

out/kutti_linux_amd64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/


cmd/kutti/rsrc_windows_amd64.syso: cmd/kutti/winres/*
	go-winres make --in=cmd/kutti/winres/winres.json --out=cmd/kutti/rsrc --arch=amd64 --product-version=${VERSION_STRING} --file-version=${VERSION_STRING}

out/kutti_windows_amd64.exe: $(KUTTICMDFILES) cmd/kutti/rsrc_windows_amd64.syso
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/completions:
	mkdir -p out/completions

out/completions/kutti-completion.ps1: $(KUTTICMDFILES) | out/completions
	go run ./cmd/kutti completion powershell > $@

out/completions/kutti.bash: $(KUTTICMDFILES) | out/completions
	go run ./cmd/kutti completion bash > $@

out/completions/_kutti: $(KUTTICMDFILES) | out/completions
	go run ./cmd/kutti completion zsh > $@

out/completions/kutti.fish: $(KUTTICMDFILES) | out/completions
	go run ./cmd/kutti completion fish > $@

.PHONY: completions
completions: out/completions/kutti.bash out/completions/_kutti out/completions/kutti.fish out/completions/kutti-completion.ps1

out/kutti-windows-installer.exe: build/package/kutti-windows-installer/kutti-windows-installer.nsi out/kutti_windows_amd64.exe out/completions/kutti-completion.ps1
	makensis -NOCD -V3 -DVERSION="${VERSION_STRING}" -DVERSION_FOUR_PART="${VERSION_FOUR_PART}" -- $<

out/kutti_$(VERSION_STRING)_linux_amd64.deb: build/package/nfpm/nfpm.yaml out/kutti_linux_amd64 manpagedocs completions
	VERSION=$(VERSION_STRING) ARCH=amd64 nfpm package -f $< -p deb -t $@

out/kutti-$(VERSION_STRING)-1.x86_64.rpm: build/package/nfpm/nfpm.yaml out/kutti_linux_amd64 manpagedocs completions
	VERSION=$(VERSION_STRING) ARCH=amd64 nfpm package -f $< -p rpm -t $@

out/kutti-$(VERSION_STRING)-1-x86_64.pkg.tar.zst: build/package/nfpm/nfpm.yaml out/kutti_linux_amd64 manpagedocs completions
	VERSION=$(VERSION_STRING) ARCH=amd64 nfpm package -f $< -p archlinux -t $@

.PHONY: deb
deb: out/kutti_$(VERSION_STRING)_linux_amd64.deb

.PHONY: rpm
rpm: out/kutti-$(VERSION_STRING)-1.x86_64.rpm

.PHONY: archlinux
archlinux: out/kutti-$(VERSION_STRING)-1-x86_64.pkg.tar.zst

.PHONY: linux-packages
linux-packages: deb rpm archlinux

out/kutti_darwin_amd64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

out/kutti_darwin_arm64: $(KUTTICMDFILES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $@ -ldflags "-X main.version=${VERSION_STRING}" ./cmd/kutti/

.PHONY: linux
linux: out/kutti_linux_amd64

.PHONY: windows
windows: out/kutti_windows_amd64.exe

.PHONY: windows-installer
windows-installer: out/kutti-windows-installer.exe

.PHONY: mac-intel
mac-intel: out/kutti_darwin_amd64

.PHONY: mac
mac: out/kutti_darwin_arm64

.PHONY: all
all: linux windows mac mac-intel

.PHONY: installers
installers: linux-packages windows-installer

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
