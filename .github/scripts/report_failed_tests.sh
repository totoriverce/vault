#!/bin/bash

set -e
MAX_TESTS=10
# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR_NUMBER:?} ]
if [ -z "$TABLE_DATA" ]; then
  BODY="CI Results:
All Go tests succeeded! :white_check_mark:"
else
  # Remove any rows that don't have a test name
  # Only keep the test type, test name, and logs column
  # Remove the scroll emoji
  TABLE_DATA=$(echo "$TABLE_DATA" | awk -F\| '{if ($4 != " - ") { print "|" $2 "|" $4 "|" $7 }}' | sed -r 's/ :scroll://')
  NUM_FAILURES=$(wc -l <<< "$TABLE_DATA")
  
  # Check if the number of failures is greater than the maximum tests to display
  # If so, limit the table to MAX_TESTS number of results
  if [ "$NUM_FAILURES" -gt "$MAX_TESTS" ]; then
      TABLE_DATA=$(echo "$TABLE_DATA" | head -n "$MAX_TESTS")
      NUM_OTHER=( $NUM_FAILURES - "$MAX_TESTS" )
      TABLE_DATA="$TABLE_DATA

and $NUM_OTHER other tests"
  fi
  
  # Add the header for the table
  BODY="CI Results:
| Test Type | Test | Logs |
| --------- | ---- | ---- |
${TABLE_DATA}"
fi

source ./.github/scripts/gh_comment.sh

update_or_create_comment "$REPO" "$PR_NUMBER" "CI Results:" "$BODY"