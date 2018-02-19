#!/bin/bash

set -euo pipefail

APPS=(
  build-essential
  # docker
  apt-transport-https
  ca-certificates
  curl
  gnupg2
  software-properties-common
)

function usage() {
  echo "$0 <docker run as user>"
  exit 1
}

[ "$#" -eq 0 ] && usage

RUN_AS=$1

export DEBIAN_FRONTEND=noninteractive

apt-get update -y
apt-get install -y ${APPS[@]}
curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") $(lsb_release -cs) stable"
apt-get update -y
apt-get install -y docker-ce
usermod -aG docker ${RUN_AS}

unset DEBIAN_FRONTEND
