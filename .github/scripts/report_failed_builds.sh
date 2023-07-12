#!/bin/bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR:?} ]
[ ${BUILD_OTHER:?} ]
[ ${BUILD_LINUX:?} ]
[ ${BUILD_DARWIN:?} ]
[ ${BUILD_DOCKER:?} ]
[ ${BUILD_UBI:?} ]
[ ${TEST:?} ]
[ ${TEST_DOCKER_K8S:?} ]

# listing out all of the jobs with the status
jobs=( "build-other:$BUILD_OTHER" "build-linux:$BUILD_LINUX" "build-darwin:$BUILD_DARWIN" "build-docker:$BUILD_DOCKER" "build-ubi:$BUILD_UBI" "test:$TEST" "test-docker-k8s:$TEST_DOCKER_K8S" )

# there is a case where even if a job is failed, it reports as cancelled. So, we look for both.
failed_jobs=()
for job in "${jobs[@]}";do
  if [[ "$job" == *"failure"* || "$job" == *"cancelled"* ]]; then
    failed_jobs+=("$job")
  fi
done

gh pr comment "$PR" --body "build failed for these jobs: ${failed_jobs[*]}. Please refer to this workflow to learn more: https://github.com/hashicorp/vault/actions/runs/$RUN_ID"
