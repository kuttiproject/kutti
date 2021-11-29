#!/bin/sh

# This script generates a POSIX sh script for downloading and installing 
# the kutti binary. The generated script is hardcoded to download a 
# specific version of kutti, for a specific OS and architecture. The
# desired version, OS and architecture have to be specified in the
# CURRENT_VERSION, GOOS and GOARCH environment variable respectively
# while calling this script.
# An environment variable called CHECKRELEASE can be set to any value to
# check if the desired version exists as a GitHub release before the 
# script is generated.

if [ -z "$CURRENT_VERSION" ];then
    echo >&2 "Please specify the kutti release version in the CURRENT_VERSION environment variable"
    exit 1
fi

if [ -z "$GOOS" ] || [ -z "$GOARCH" ]; then
    echo >&2 "Please specify the target OS and architecture in the GOOS and GOARCH environment variables."
    exit 1
fi

if [ -n "$CHECKRELEASE" ]; then
    CHECKCMD=""
    if [ -n "$(which curl)" ]; then
        CHECKCMD="curl -L -O /dev/null --head --silent --fail https://github.com/kuttiproject/kutti/releases/download/v${CURRENT_VERSION}/kutti_${GOOS}_${GOARCH}"
    elif [ -n "$(which wget)" ]; then
        CHECKCMD="wget -q -O /dev/null --spider https://github.com/kuttiproject/kutti/releases/download/v${CURRENT_VERSION}/kutti_${GOOS}_${GOARCH}"
    else
        echo >&2 "Need either curl or wget to generate install script."
        exit 1
    fi

    if [ ! "$($CHECKCMD)" ]; then
        echo >&2 "The combination Version=${CURRENT_VERSION}, OS:$GOOS and Arch:$GOARCH is not currently available."
        exit 2
    fi
fi

cat <<EOSCRIPT
#!/bin/sh

download() {
    echo "Downloading kutti version ${CURRENT_VERSION} for $GOOS/$GOARCH"

    if [ -n "\$(which curl)" ]; then
        curl -LO https://github.com/kuttiproject/kutti/releases/download/v${CURRENT_VERSION}/kutti_${GOOS}_${GOARCH}
    elif [ -n "\$(which wget)" ]; then
        echo "wget -q -O kutti_${GOOS}_${GOARCH} https://github.com/kuttiproject/kutti/releases/download/v${CURRENT_VERSION}/kutti_${GOOS}_${GOARCH}"
    else
        echo >&2 "Need either curl or wget to install kutti."
        exit 1
    fi

    chmod +x kutti_${GOOS}_${GOARCH}

    echo "Copying kutti to /usr/local/bin"
    sudo mv kutti_${GOOS}_${GOARCH} /usr/local/bin/kutti
    
    echo "Done. Installed version of kutti is:"
    kutti -v
}

download
EOSCRIPT
