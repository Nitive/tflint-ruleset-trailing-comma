# TFLint Ruleset Trailing Comma
[![Build Status](https://github.com/Nitive/tflint-ruleset-trailing-comma/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/Nitive/tflint-ruleset-trailing-comma/actions)

This plugin enforces a compact comma style for Terraform expressions:

- Multiline lists and multiline function calls must end with a trailing comma.
- Single-line lists and single-line function calls must not end with a trailing comma.
- Multiline maps must not contain commas.
- Single-line maps may use commas between entries, but must not end with a trailing comma.

See also [Writing Plugins](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/plugins.md).

## Requirements

- TFLint v0.46+
- Go v1.26

## Installation

This repository publishes a local/custom TFLint plugin. Build it locally for development, or use GitHub releases once you wire publishing for your repo.

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "tflint-ruleset-trailing-comma" {
  enabled = true

  version = "0.1.0"
  source  = "github.com/Nitive/tflint-ruleset-trailing-comma"
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|multiline_trailing_comma|Enforce trailing commas for multiline lists and function calls, while forbidding them on single-line forms|ERROR|✔||
|multiline_map_no_comma|Forbid commas in multiline maps and trailing commas in single-line maps|ERROR|✔||

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

You can run the built plugin like the following:

```
$ cat << EOS > .tflint.hcl
plugin "tflint-ruleset-trailing-comma" {
  enabled = true
}
EOS
$ tflint
```
