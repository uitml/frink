#!/bin/sh

set -eu
umask 077

api_url="https://raw.githubusercontent.com/uitml/frink/master"

mkdir -p /tmp/frink && cd /tmp/frink
curl -LO "${api_url}/bin/kubectl-job-kill"
curl -LO "${api_url}/bin/kubectl-job-list"
curl -LO "${api_url}/bin/kubectl-job-run"
curl -LO "${api_url}/bin/kubectl-job-watch"
sudo cp kubectl* /usr/local/bin/
sudo chmod 755 /usr/local/bin/kubectl-job-*
