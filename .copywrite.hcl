schema_version = 1

project {
  license        = "MPL-2.0"
  copyright_year = 2015

  # (OPTIONAL) A list of globs that should not have copyright/license headers.
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  header_ignore = [
    "builtin/credential/aws/pkcs7/**",
    "ui/node_modules/**",
    "enos/modules/k8s_deploy_vault/raft-config.hcl",
    "plugins/database/postgresql/scram/**"
  ]
}
