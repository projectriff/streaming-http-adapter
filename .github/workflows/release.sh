#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# from https://cloud.google.com/sdk/docs/downloads-apt-get
export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)"
echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo apt-get update && sudo apt-get install google-cloud-sdk

gcloud config set disable_prompts True
gcloud auth activate-service-account --key-file <(echo ${GCLOUD_CLIENT_SECRET} | base64 --decode)

readonly root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." >/dev/null 2>&1 && pwd)
readonly version=$(cat ${root}/VERSION)
readonly git_branch=${GITHUB_REF:11} # drop 'refs/head/' prefix
readonly git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
readonly slug=${version}-${git_timestamp}-${GITHUB_SHA:0:16}

echo "Publishing riff http->streaming adapter"
(cd $root && gsutil cp -a public-read streaming-http-adapter-linux-amd64.tgz gs://projectriff/streaming-http-adapter/streaming-http-adapter-linux-amd64-${slug}.tgz)
