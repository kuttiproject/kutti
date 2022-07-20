Param(
    $VersionMajor = (property VERSION_MAJOR "0"),
    $VersionMinor = (property VERSION_MINOR "3"),
    $BuildNumber  = (property BUILD_NUMBER "2"),
    $PatchString  = (property PATCH_NUMBER  "")
)

# Maintain semantic version in the parameters above
# Also change in cmd/kutti/main.go
$VersionString = "$($VersionMajor).$($VersionMinor).$($BuildNumber)$($PatchString)"

$KuttiCmdFiles = (Get-Item "cmd/kutti/*.go") +          `
				 (Get-Item "internal/pkg/cli/*.go") +   `
				 (Get-Item "internal/pkg/cmd/*.go") +   `
				 (Get-Item "internal/pkg/cmd/*/*.go") + `
				 (Get-Item "go.mod") +                  `
                 (Get-Item "invoke.build.ps1") 

# Synopsis: Show usage
task . {
	Write-Host "Usage: make linux|windows|mac|linux-install-script|windows-installer|mac-install-script|all|installers|clean"
}

# Synopsis: Build output directory
task outputdir -Outputs out\ {
    New-Item -Path out\ -ItemType Directory -ErrorAction Ignore
}

# Synopsis: Build linux binary
task linux -Outputs out/kutti_linux_amd64 -Inputs $($KuttiCmdFiles) {
    exec {
        $env:CGO_ENABLED="0"
        $env:GOOS="linux"
        $env:GOARCH="amd64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
    }
}

# Synopsis: Build windows resource file
task winres -Outputs cmd/kutti/rsrc_windows_amd64.syso -Inputs (Get-Item "cmd/kutti/winres/*") {
    exec {
        go-winres make --in=cmd/kutti/winres/winres.json --out=cmd/kutti/rsrc --arch=amd64 --product-version=$($VersionString) --file-version=$($VersionString)
    }
}

# Synopsis: Build windows binary
task windows -Outputs out/kutti_windows_amd64.exe -Inputs {$($KuttiCmdFiles) + (Get-Item -Path "cmd/kutti/rsrc_windows_amd64.syso")} winres, {
    exec {
        $env:CGO_ENABLED="0"
        $env:GOOS="windows"
        $env:GOARCH="amd64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
    }
}

# Synopsis: Build mac binary
task mac -Outputs out/kutti_darwin_amd64 -Inputs $($KuttiCmdFiles) {
    exec {
        $env:CGO_ENABLED="0"
        $env:GOOS="darwin"
        $env:GOARCH="amd64"
        go build -o $($Outputs) -ldflags "-X main.version=$($VersionString)" ./cmd/kutti/
    }
}

# Synopsis: Build linux installation script
task linux-install-script -Outputs out/get-kutti-linux.sh -Inputs build/package/posix-install-script/generate-script.ps1 outputdir, {
    $env:CURRENT_VERSION=$VersionString
    $env:GOOS="linux"
    $env:GOARCH="amd64"

    exec {
        Invoke-Expression  $Inputs[0] > $Outputs
    }
}

# Synopsis: Build windows installer
task windows-installer -Outputs out/kutti-windows-installer.exe -Inputs build/package/kutti-windows-installer/kutti-windows-installer.nsi windows, {
	makensis -NOCD -V3 -- $($Inputs[0])
}

# Synopsis: Build mac installation script
task mac-install-script -Outputs out/get-kutti-darwin.sh -Inputs build/package/posix-install-script/generate-script.ps1 outputdir, {
    $env:CURRENT_VERSION=$VersionString
    $env:GOOS="darwin"
    $env:GOARCH="amd64"

    exec {
        Invoke-Expression  $Inputs[0] > $Outputs
    }
}

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
task all linux, windows, mac

# Synopsis: Build all installers
task installers linux-install-script, windows-installer, mac-install-script

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
