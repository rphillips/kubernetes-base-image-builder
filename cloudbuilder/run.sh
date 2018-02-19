#!/bin/bash

set -euo pipefail

export REGISTRY=${REGISTRY:-gcr.io/rphillips-dev}

WORKSPACE=$(dirname "$0")
cd $WORKSPACE

archs=(
  amd64
  arm
  arm64
  ppc64le
  s390x
)

function clone() {
  [ -d kubernetes ] && return
  git clone --progress https://github.com/kubernetes/kubernetes
}

function _build() {
  export ARCH=$1 TAG=$(date +%Y%m%d)
  pushd kubernetes/build/debian-base
  make build
  docker save $REGISTRY/debian-base-$ARCH:$TAG | gzip > ${WORKSPACE}/debian-base-$ARCH.tar.gz
  unset ARCH TAG
  popd
}

function build() {
  for arch in ${archs[@]}; do
    _build $arch
  done
}

function cleanup() {
  rm -rf kubernetes
}

clone
build
cleanup
