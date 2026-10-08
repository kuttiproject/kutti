Param(
    $VersionMajor = (property VERSION_MAJOR "0"),
    $VersionMinor = (property VERSION_MINOR "6"),
    $BuildNumber = (property BUILD_NUMBER "0"),
    $PatchString = (property PATCH_NUMBER  "")
)

# Maintain semantic version in the parameters above
# Also change in cmd/kutti/main.go
$VersionString = "$($VersionMajor).$($VersionMinor).$($BuildNumber)$($PatchString)"
$VersionFourPart = "$($VersionMajor).$($VersionMinor).$($BuildNumber).0"

$KuttiCmdFiles = (Get-Item "cmd/kutti/*.go") + `
(Get-Item "internal/pkg/cli/*.go") + `
(Get-Item "internal/pkg/cmd/*.go") + `
(Get-Item "internal/pkg/cmd/*/*.go") + `
(Get-Item "go.mod") + `
(Get-Item "invoke.build.ps1") 

# Synopsis: Show usage
task . {
    Write-Host "Usage: Invoke-Build linux|windows|mac|mac-intel|deb|rpm|archlinux|linux-packages|completions|windows-installer|all|installers|clean"
}

# Synopsis: Build output directory
task outputdir -Outputs out\ {
    New-Item -Path out\ -ItemType Directory -ErrorAction Ignore
}

# Synopsis: Build linux binary
task linux -Outputs out/kutti_linux_amd64 -Inputs $($KuttiCmdFiles) {
    exec {
        $env:CGO_ENABLED = "0"
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
        Remove-Item env:\CGO_ENABLED -ErrorAction SilentlyContinue
        Remove-Item env:\GOOS -ErrorAction SilentlyContinue
        Remove-Item env:\GOARCH -ErrorAction SilentlyContinue
    }
}

# Synopsis: Build windows resource file
task winres -Outputs cmd/kutti/rsrc_windows_amd64.syso -Inputs (Get-Item "cmd/kutti/winres/*") {
    exec {
        go-winres make --in=cmd/kutti/winres/winres.json --out=cmd/kutti/rsrc --arch=amd64 --product-version=$($VersionString) --file-version=$($VersionString)
    }
}

# Synopsis: Build windows binary
task windows -Outputs out/kutti_windows_amd64.exe -Inputs { $($KuttiCmdFiles) + (Get-Item -Path "cmd/kutti/rsrc_windows_amd64.syso") } winres, {
    exec {
        $env:CGO_ENABLED = "0"
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
        Remove-Item env:\CGO_ENABLED -ErrorAction SilentlyContinue
        Remove-Item env:\GOOS -ErrorAction SilentlyContinue
        Remove-Item env:\GOARCH -ErrorAction SilentlyContinue
    }
}

# Synopsis: Build completions output directory
task completionsoutputdir -Outputs out\completions outputdir, {
    New-Item out\completions -ItemType Directory -ErrorAction Ignore
}

# Synopsis: Build PowerShell completions
task powershell-completion -Outputs out\completions\kutti-completion.ps1 -Inputs $($KuttiCmdFiles) completionsoutputdir, {
    exec {
        Write-Output "GOOS: $($env:GOOS)"
        go run ./cmd/kutti completion powershell > $($Outputs)
    }
}

# Synopsis: Build Bash completions
task bash-completion -Outputs out/completions/kutti.bash -Inputs $($KuttiCmdFiles) completionsoutputdir, {
    exec {
        go run ./cmd/kutti completion bash > $($Outputs)
    }
}

# Synopsis: Build Zsh completions
task zsh-completion -Outputs out/completions/_kutti -Inputs $($KuttiCmdFiles) completionsoutputdir, {
    exec {
        go run ./cmd/kutti completion zsh > $($Outputs)
    }
}

# Synopsis: Build Fish completions
task fish-completion -Outputs out/completions/kutti.fish -Inputs $($KuttiCmdFiles) completionsoutputdir, {
    exec {
        go run ./cmd/kutti completion fish > $($Outputs)
    }
}

# Synopsis: Build all shell completions
task completions bash-completion, zsh-completion, fish-completion, powershell-completion

# Synopsis: Build mac amd64 binary
task mac-intel -Outputs out/kutti_darwin_amd64 -Inputs $($KuttiCmdFiles) {
    exec {
        $env:CGO_ENABLED = "0"
        $env:GOOS = "darwin"
        $env:GOARCH = "amd64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
        Remove-Item env:\CGO_ENABLED -ErrorAction SilentlyContinue
        Remove-Item env:\GOOS -ErrorAction SilentlyContinue
        Remove-Item env:\GOARCH -ErrorAction SilentlyContinue
    }
}

# Synopsis: Build mac arm64 binary
task mac -Outputs out/kutti_darwin_arm64 -Inputs $($KuttiCmdFiles) {
    exec {
        $env:CGO_ENABLED = "0"
        $env:GOOS = "darwin"
        $env:GOARCH = "arm64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
        Remove-Item env:\CGO_ENABLED -ErrorAction SilentlyContinue
        Remove-Item env:\GOOS -ErrorAction SilentlyContinue
        Remove-Item env:\GOARCH -ErrorAction SilentlyContinue
    }
}

# Synopsis: Build windows installer
task windows-installer -Outputs out/kutti-windows-installer.exe -Inputs build/package/kutti-windows-installer/kutti-windows-installer.nsi windows, powershell-completion, {
    makensis -NOCD -V3 "-DVERSION=$($VersionString)" "-DVERSION_FOUR_PART=$($VersionFourPart)" -- $($Inputs[0])
}

# Synopsis: Build Linux Debian package
task deb -Outputs "out/kutti_$($VersionString)_linux_amd64.deb" -Inputs build/package/nfpm/nfpm.yaml linux, manpagedocs, completions, {
    $env:VERSION = $VersionString
    $env:ARCH = "amd64"
    exec {
        nfpm package -f $Inputs[0] -p deb -t $Outputs
    }
}

# Synopsis: Build Linux RPM package
task rpm -Outputs "out/kutti-$($VersionString)-1.x86_64.rpm" -Inputs build/package/nfpm/nfpm.yaml linux, manpagedocs, completions, {
    $env:VERSION = $VersionString
    $env:ARCH = "amd64"
    exec {
        nfpm package -f $Inputs[0] -p rpm -t $Outputs
    }
}

# Synopsis: Build Arch Linux package
task archlinux -Outputs "out/kutti-$($VersionString)-1-x86_64.pkg.tar.zst" -Inputs build/package/nfpm/nfpm.yaml linux, manpagedocs, completions, {
    $env:VERSION = $VersionString
    $env:ARCH = "amd64"
    exec {
        nfpm package -f $Inputs[0] -p archlinux -t $Outputs
    }
}

# Synopsis: Build all Linux packages
task linux-packages deb, rpm, archlinux

# Synopsis: Build manpage docs output directory
task manpagedocsoutputdir -Outputs out\man {
    New-Item out\man -ItemType Directory -ErrorAction Ignore
}

# Synopsis: Build manpage docs
task manpagedocs -Outputs out/man/kutti.1 -Inputs $($KuttiCmdFiles) manpagedocsoutputdir, {
    exec {
        go run internal/cmd/gendoc/main.go -o out/man -t manpages
    }
}

# Synopsis: Build markdown docs output directory
task markdowndocsoutputdir -Outputs out\markdown {
    New-Item out\markdown -ItemType Directory -ErrorAction Ignore
}

# Synopsis: Build markdown docs
task markdowndocs -Outputs out/markdown/kutti.md -Inputs $($KuttiCmdFiles) markdowndocsoutputdir, {
    exec {
        go run internal/cmd/gendoc/main.go -o out/markdown -t markdown
    }
}

# Synopsis: Build all binaries
task all linux, windows, mac, mac-intel

# Synopsis: Build all installers
task installers linux-packages, windows-installer

# Synopsis: Build all docs
task docs manpagedocs, markdowndocs

# Synopsis: Clean built windows resource file
task resourceclean {
    Remove-Item -Force -ErrorAction Ignore ./cmd/kutti/rsrc_windows_amd64.syso
}

# Synopsis: Clean built binaries
task binclean {
    Remove-Item -Recurse -Force -ErrorAction Ignore ./out
}

# Synopsis: Clean built manpage docs
task manpagedocsclean {
    exec {
        Remove-Item -Recurse -Force -ErrorAction Ignore ./out/man
    }
}

# Synopsis: Clean built markdown docs
task markdowndocsclean {
    exec {
        Remove-Item -Recurse -Force -ErrorAction Ignore ./out/markdown
    }
}

# Synopsis: Clean all docs
task docsclean manpagedocsclean, markdowndocsclean

# Synopsis: Clean everything
task clean resourceclean, binclean
