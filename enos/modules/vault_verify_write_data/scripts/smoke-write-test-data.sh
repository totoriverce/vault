#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

function fail {
	echo "$1" 1>&2
	exit 1
}

binpath=${VAULT_INSTALL_DIR}/vault
testkey=${TEST_KEY}
testvalue=${TEST_VALUE}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

retry 5 $binpath kv put secret/test $testkey=$testvalue
