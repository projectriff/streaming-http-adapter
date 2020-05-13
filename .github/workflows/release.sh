#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

gcloud config set disable_prompts True
gcloud auth activate-service-account --key-file <(echo ${GCLOUD_CLIENT_SECRET} | base64 --decode)

readonly version=$(cat VERSION)
readonly git_branch=${GITHUB_REF:11} # drop 'refs/head/' prefix
readonly git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
readonly slug=${version}-${git_timestamp}-${GITHUB_SHA:0:16}

echo "Publishing streaming-http-adapter-linux-amd64-${slug}.tgz"
gsutil cp streaming-http-adapter-linux-amd64.tgz gs://projectriff/streaming-http-adapter/streaming-http-adapter-linux-amd64-${slug}.tgz
