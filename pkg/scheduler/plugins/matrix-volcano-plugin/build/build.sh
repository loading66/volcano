#!/bin/bash
# Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

set -e
TOP_DIR=${GOPATH}/src/volcano.sh/volcano/
BASE_PATH=${GOPATH}/src/volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/

function clean() {
    rm -f "${BASE_PATH}"/output/*.so
}

function build() {
    echo "Start Build ......"

    export PATH=$GOPATH/bin:$PATH
    export GO111MODULE=on

    cd "${TOP_DIR}"
    go mod tidy

    mkdir -p "${BASE_PATH}"/output/
    cd "${BASE_PATH}"/output/

    REL_MATRIX_PLUGIN=matrix_volcano_plugin

    CGO_CFLAGS="-fstack-protector-strong -D_FORTIFY_SOURCE=2 -O2 -fPIC -ftrapv" \
    CGO_CPPFLAGS="-fstack-protector-strong -D_FORTIFY_SOURCE=2 -O2 -fPIC -ftrapv" \
    CC=/opt/buildtools/musl-1.2.5/bin/musl-gcc CGO_ENABLED=1 \
    go build -buildvcs=false -mod=mod -buildmode=plugin -ldflags "-s -linkmode=external -extldflags=-Wl,-z,now
    -X volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin.PluginName=${REL_MATRIX_PLUGIN}" \
    -o "${REL_MATRIX_PLUGIN}".so "${GOPATH}"/src/volcano.sh/volcano/pkg/scheduler/plugins/matrix-volcano-plugin/

    if [ ! -f "${BASE_PATH}/output/${REL_MATRIX_PLUGIN}.so" ]
    then
      echo "fail to find ${REL_MATRIX_PLUGIN}.so"
      exit 1
    fi

    chmod 400 "${BASE_PATH}"/output/*.so
}

function main() {
  clean
  build
}

main "${1}"

echo ""
echo "Build Finished!"
echo ""