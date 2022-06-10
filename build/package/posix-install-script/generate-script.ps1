# This script generates a POSIX sh script for downloading and installing 
# the kutti binary. The generated script is hardcoded to download a 
# specific version of kutti, for a specific OS and architecture. The
# desired version, OS and architecture have to be specified in the
# CURRENT_VERSION, GOOS and GOARCH environment variable respectively
# while calling this script.
# An environment variable called CHECKRELEASE can be set to any value to
# check if the desired version exists as a GitHub release before the 
# script is generated.

If ( [String]::IsNullOrEmpty($env:CURRENT_VERSION)) {
    $Host.UI.WriteErrorLine("Please specify the kutti release version in the CURRENT_VERSION environment variable")
    exit 1
}

If ([String]::IsNullOrEmpty($env:GOOS) -or [String]::IsNullOrEmpty($env:GOARCH) ){
    $Host.UI.WriteErrorLine( "Please specify the target OS and architecture in the GOOS and GOARCH environment variables.")
    exit 1
}

If  (![String]::IsNullOrEmpty($env:CHECKRELEASE)) {

    Invoke-WebRequest "https://github.com/kuttiproject/kutti/releases/download/v$($env:CURRENT_VERSION)/kutti_$($env:GOOS)_$($env:GOARCH)" -Method Head -ErrorAction Continue -ErrorVariable ev
    If (![String]::IsNullOrEmpty($ev)) {
        $Host.UI.WriteErrorLine("The combination Version=$($env:CURRENT_VERSION), OS:$($env:GOOS) and Arch:$($env:GOARCH) is not currently available.")
        exit 2        
    }
}

Write-Output @"
#!/bin/sh

download() {
    echo "Downloading kutti version $($env:CURRENT_VERSION) for $($env:GOOS)/$($env:GOARCH)"

    if [ -n "`$(which curl)" ]; then
        curl -LO https://github.com/kuttiproject/kutti/releases/download/v$($CURRENT_VERSION)/kutti_$($env:GOOS)_$($env:GOARCH)
    elif [ -n "`$(which wget)" ]; then
        echo "wget -q -O kutti_$($env:GOOS)_$($env:GOARCH) https://github.com/kuttiproject/kutti/releases/download/v$($env:CURRENT_VERSION)/kutti_$($env:GOOS)_$($env:GOARCH)"
    else
        echo >&2 "Need either curl or wget to install kutti."
        exit 1
    fi

    chmod +x kutti_$($env:GOOS)_$($env:GOARCH)

    echo "Copying kutti to /usr/local/bin"
    sudo mv kutti_$($env:GOOS)_$($env:GOARCH) /usr/local/bin/kutti
    
    echo "Done. Installed version of kutti is:"
    kutti -v
}

download
"@
