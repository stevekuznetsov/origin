#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GO_VERSION=($(go version))

if [[ -z $(echo "${GO_VERSION[2]}" | grep -E 'go1.4') ]]; then
  echo "Unknown go version '${GO_VERSION}', skipping go vet."
  exit 0
fi

OS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${OS_ROOT}/hack/common.sh"
source "${OS_ROOT}/hack/util.sh"

cd "${OS_ROOT}"

dirs=$(find_files | cut --delimiter=/ --fields=1-2 | sort -u)
for dir in $dirs
do
  go tool vet -all=true -shadow=false -composites=false -copylocks=false ${dir} 2>&1
done
