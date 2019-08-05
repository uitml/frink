#!/bin/sh

set -eu
umask 077

api_url="https://raw.githubusercontent.com/uitml/frink/master"

mkdir -p /tmp/frink && cd /tmp/frink
curl -LO "${api_url}/bin/kubectl-job-kill"
curl -LO "${api_url}/bin/kubectl-job-list"
curl -LO "${api_url}/bin/kubectl-job-run"
curl -LO "${api_url}/bin/kubectl-job-watch"
chmod +x kubectl*
sudo cp --no-preserve=ownership kubectl* /usr/local/bin/
