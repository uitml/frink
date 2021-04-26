#!/bin/sh

set -eu
umask 077

kernel="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${kernel}" in
  linux*) os="linux" ;;
  darwin*) os="macos" ;;
  *) echo "unknown kernel" && exit 1 ;;
esac

arch="$(uname -m)"
case "${arch}" in
  x86_64*) arch="amd64" ;;
esac

api_url="https://raw.githubusercontent.com/uitml/frink/legacy"
version="$(curl -s https://api.github.com/repos/uitml/frink/releases/latest | grep "\"name\": \"v.*\"" | cut -d '"' -f 4)"
download_url="https://github.com/uitml/frink/releases/download/${version}/frink-${os}-${arch}"

echo ${download_url}

cd "$(mktemp -d -t frink-XXXX)"
curl -LO "${download_url}"
sudo cp frink* /usr/local/bin/frink
sudo chmod 755 /usr/local/bin/frink
