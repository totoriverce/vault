# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This script is sourced into the shell running in a Github Actions workflow.

# Usage:
# To check deprecations locally using the script, follow these steps:
# From the repository root or within a package folder, execute deprecations-checker.sh 
# Optionally: to only show deprecations in changed files between the current branch and 
# a specific branch, pass the other branch name as an argument to the script.
#
# For example: 
# ./scripts/deprecations-checker.sh (or) make deprecations
# ./scripts/deprecations-checker.sh main (or) make ci-deprecations
#
# If no branch name is specified, the command will show all usage of deprecations in the code.
#
# GitHub Actions runs this against the PR's base ref branch. 

# Staticcheck uses static analysis to finds bugs and performance issues, offers simplifications, 
# and enforces style rules.
# Here, it is used to check if a deprecated function, variable, constant or field is used.

# Run staticcheck 
echo "Performing deprecations check: running staticcheck"

# Identify repository name
if [ -z $2 ]; then
    # local repository name
    repositoryName=$(basename `git rev-parse --show-toplevel`)
else
    # github repository name from deprecated-functions-checker.yml
    repositoryName=$2
fi

# Modify the command with the correct build tag based on repository
if [ $repositoryName == "vault-enterprise" ]; then
    staticcheckCommand=$(echo "staticcheck ./... -tags=enterprise")
else
    staticcheckCommand=$(echo "staticcheck ./...")
fi

# If no compare branch name is specified, output all deprecations
# Else only output the deprecations from the changes added
if [ -z $1 ]
    then
        $staticcheckCommand | grep deprecated
    else
        # GitHub Actions will use this to find only changes wrt PR's base ref branch
        # revgrep CLI tool will return an exit status of 1 if any issues match, else it will return 0
        $staticcheckCommand | grep deprecated 2>&1 | revgrep "$(git merge-base HEAD "origin/$1")"
fi
