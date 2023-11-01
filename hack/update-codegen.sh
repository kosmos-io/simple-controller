#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CLUSTERLINK_GROUP_PACKAGE="github.com"
CLUSTERLINK_GO_PACKAGE="${CLUSTERLINK_GROUP_PACKAGE}/kosmos.io/simple-controller"

# For all commands, the working directory is the parent directory(repo root).
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}"

echo "Generating with deepcopy-gen"
GO111MODULE=on go install k8s.io/code-generator/cmd/deepcopy-gen
export GOPATH=$(go env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

group_path="${REPO_ROOT}/${CLUSTERLINK_GROUP_PACKAGE}"
link_path="${REPO_ROOT}/${CLUSTERLINK_GO_PACKAGE}"
function cleanup() {
  rm -rf "${group_path}"
}
trap "cleanup" EXIT SIGINT
cleanup

function create_gopath_tree() {
  local repo_root=$1
  local link_path=$2

  go_pkg_dir=$(dirname "${link_path}")

  mkdir -p "${go_pkg_dir}"

  if [[ ! -e "${go_pkg_dir}" || "$(readlink "${go_pkg_dir}")" != "${repo_root}" ]]; then
    ln -snf "${repo_root}" "${go_pkg_dir}"
  fi
}
create_gopath_tree "${REPO_ROOT}" "${link_path}"

deepcopy-gen \
  --go-header-file hack/boilerplate.go.txt \
  --input-dirs="github.com/kosmos.io/simple-controller/pkg/apis/v1" \
  --output-base="${REPO_ROOT}" \
  --output-package="pkg/apis/v1" \
  --output-file-base=zz_generated.deepcopy

echo "Generating with register-gen"
GO111MODULE=on go install k8s.io/code-generator/cmd/register-gen
register-gen \
  --go-header-file hack/boilerplate.go.txt \
  --input-dirs="github.com/kosmos.io/simple-controller/pkg/apis/v1" \
  --output-base="${REPO_ROOT}" \
  --output-package="pkg/apis/v1" \
  --output-file-base=zz_generated.register