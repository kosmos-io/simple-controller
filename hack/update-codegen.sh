#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GROUP_PACKAGE="github.com"
GO_PACKAGE="${GROUP_PACKAGE}/kosmos.io/simple-controller"

# For all commands, the working directory is the parent directory(repo root).
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}"

echo "Generating with deepcopy-gen"
GO111MODULE=on go install k8s.io/code-generator/cmd/deepcopy-gen
export GOPATH=$(go env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

function cleanup() {
  rm -rf "${REPO_ROOT}/${GROUP_PACKAGE}"
}
trap "cleanup" EXIT SIGINT
cleanup

mkdir -p "$(dirname "${REPO_ROOT}/${GO_PACKAGE}")"

deepcopy-gen \
  --input-dirs="github.com/kosmos.io/simple-controller/pkg/apis/v1" \
  --output-base="${REPO_ROOT}" \
  --output-package="pkg/apis/v1" \
  --output-file-base=zz_generated.deepcopy

echo "Generating with register-gen"
GO111MODULE=on go install k8s.io/code-generator/cmd/register-gen
register-gen \
  --input-dirs="github.com/kosmos.io/simple-controller/pkg/apis/v1" \
  --output-base="${REPO_ROOT}" \
  --output-package="pkg/apis/v1" \
  --output-file-base=zz_generated.register

mv "${REPO_ROOT}/${GO_PACKAGE}"/pkg/apis/v1/* "${REPO_ROOT}"/pkg/apis/v1