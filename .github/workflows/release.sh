#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." >/dev/null 2>&1 && pwd)
readonly version=$(cat ${root}/VERSION)
readonly git_branch=${GITHUB_REF:11} # drop 'refs/head/' prefix
readonly git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
readonly slug=${version}-${git_timestamp}-${GITHUB_SHA:0:16}

echo "Publishing riff http->streaming adapter"
(cd $root && gsutil cp -a public-read streaming-http-adapter-linux-amd64.tgz gs://projectriff/streaming-http-adapter-buildpack/streaming-http-adapter-linux-amd64-${slug}.tgz)
