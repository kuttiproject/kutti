#!/bin/sh
# get-kutti.sh - Universal installer for the kutti CLI
#
# Usage:
#   curl -sSL https://github.com/kuttiproject/kutti/releases/latest/download/get-kutti.sh | sh
#
# Environment variables:
#   KUTTI_VERSION     Version to install (default: latest)
#   KUTTI_INSTALL_DIR Install directory (default: ~/.local/bin, or /usr/local/bin if root)
#   KUTTI_NO_VERIFY   Set to 1 to skip SHA256 checksum verification

set -e

REPO="kuttiproject/kutti"
GITHUB_BASE="https://github.com/${REPO}/releases"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

say() {
    printf 'get-kutti: %s\n' "$1"
}

err() {
    printf 'get-kutti: error: %s\n' "$1" >&2
    exit 1
}

# ---------------------------------------------------------------------------
# Platform detection
# ---------------------------------------------------------------------------

detect_os() {
    _os=$(uname -s)
    case "${_os}" in
        Linux)  echo "linux" ;;
        Darwin) echo "darwin" ;;
        *)      err "Unsupported OS: ${_os}. Only Linux and macOS are supported." ;;
    esac
}

detect_arch() {
    _arch=$(uname -m)
    case "${_arch}" in
        x86_64)          echo "amd64" ;;
        aarch64 | arm64) echo "arm64" ;;
        *)               err "Unsupported architecture: ${_arch}. Only amd64 and arm64 are supported." ;;
    esac
}

# ---------------------------------------------------------------------------
# Download helpers (curl or wget)
# ---------------------------------------------------------------------------

_downloader=""
detect_downloader() {
    if command -v curl > /dev/null 2>&1; then
        _downloader="curl"
    elif command -v wget > /dev/null 2>&1; then
        _downloader="wget"
    else
        err "Neither curl nor wget found. Please install one and try again."
    fi
}

download() {
    _url="$1"
    _dest="$2"
    case "${_downloader}" in
        curl) curl -sSL --fail -o "${_dest}" "${_url}" ;;
        wget) wget -qO "${_dest}" "${_url}" ;;
    esac
}

# ---------------------------------------------------------------------------
# Resolve latest version via GitHub redirect (no API token or jq needed)
# ---------------------------------------------------------------------------

get_latest_version() {
    case "${_downloader}" in
        curl)
            curl -sSL -o /dev/null -w '%{url_effective}' "${GITHUB_BASE}/latest" \
                | sed 's|.*/v||'
            ;;
        wget)
            wget -qSO /dev/null "${GITHUB_BASE}/latest" 2>&1 \
                | grep -i "Location:" \
                | tail -1 \
                | sed 's|.*/v||' \
                | tr -d '\r\n '
            ;;
    esac
}

# ---------------------------------------------------------------------------
# Checksum verification
# ---------------------------------------------------------------------------

verify_checksum() {
    _binary_file="$1"
    _binary_name="$2"
    _checksum_file="$3"

    _expected=$(grep "${_binary_name}" "${_checksum_file}" | awk '{print $1}')
    if [ -z "${_expected}" ]; then
        err "Could not find checksum for '${_binary_name}' in checksums.txt"
    fi

    if command -v sha256sum > /dev/null 2>&1; then
        _actual=$(sha256sum "${_binary_file}" | awk '{print $1}')
    elif command -v shasum > /dev/null 2>&1; then
        _actual=$(shasum -a 256 "${_binary_file}" | awk '{print $1}')
    else
        say "Warning: no SHA256 tool found; skipping checksum verification."
        return 0
    fi

    if [ "${_actual}" != "${_expected}" ]; then
        err "Checksum mismatch for '${_binary_name}'!
  Expected: ${_expected}
  Got:      ${_actual}"
    fi

    say "Checksum verified."
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

# Detect platform
GOOS=$(detect_os)
GOARCH=$(detect_arch)
say "Detected platform: ${GOOS}/${GOARCH}"

# Detect download tool
detect_downloader

# Resolve version
if [ -z "${KUTTI_VERSION}" ]; then
    say "Determining latest version..."
    KUTTI_VERSION=$(get_latest_version)
    if [ -z "${KUTTI_VERSION}" ]; then
        err "Could not determine latest version. Set KUTTI_VERSION explicitly and try again."
    fi
fi
say "Version: ${KUTTI_VERSION}"

# Determine install directory
if [ -z "${KUTTI_INSTALL_DIR}" ]; then
    if [ "$(id -u)" -eq 0 ]; then
        KUTTI_INSTALL_DIR="/usr/local/bin"
    else
        KUTTI_INSTALL_DIR="${HOME}/.local/bin"
    fi
fi
say "Install directory: ${KUTTI_INSTALL_DIR}"

# Construct names and URLs
BINARY_NAME="kutti_${GOOS}_${GOARCH}"
DOWNLOAD_BASE="${GITHUB_BASE}/download/v${KUTTI_VERSION}"
BINARY_URL="${DOWNLOAD_BASE}/${BINARY_NAME}"
CHECKSUM_URL="${DOWNLOAD_BASE}/checksums.txt"

# Create temp working directory, cleaned up on exit
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT INT TERM

# Download binary
say "Downloading ${BINARY_NAME}..."
download "${BINARY_URL}" "${TMP_DIR}/${BINARY_NAME}"

# Verify checksum
if [ "${KUTTI_NO_VERIFY}" = "1" ]; then
    say "Skipping checksum verification (KUTTI_NO_VERIFY=1)."
else
    say "Downloading checksums.txt..."
    download "${CHECKSUM_URL}" "${TMP_DIR}/checksums.txt"
    say "Verifying checksum..."
    verify_checksum "${TMP_DIR}/${BINARY_NAME}" "${BINARY_NAME}" "${TMP_DIR}/checksums.txt"
fi

# Install
mkdir -p "${KUTTI_INSTALL_DIR}"
chmod 755 "${TMP_DIR}/${BINARY_NAME}"
cp "${TMP_DIR}/${BINARY_NAME}" "${KUTTI_INSTALL_DIR}/kutti"
say "Installed: ${KUTTI_INSTALL_DIR}/kutti"

# PATH advisory
case ":${PATH}:" in
    *":${KUTTI_INSTALL_DIR}:"*)
        say "Done. Run 'kutti version' to verify the installation."
        ;;
    *)
        say "Done."
        say ""
        say "Note: ${KUTTI_INSTALL_DIR} is not in your PATH."
        say "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        say "  export PATH=\"${KUTTI_INSTALL_DIR}:\$PATH\""
        ;;
esac
