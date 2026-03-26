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

  signing_key = <<-KEY
-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGnFFa4BEAC/srSOdGEAAipn0SOWg0de1MjOOf8FZ0k9kNLMPDCo3AQ4gtLy
4/jWmX3oiZG8Ial2F/9wKnypk79pBxpR8y7mXLMbg783zSk2A82SEZtnjMcG/EAr
X4tRzG5lrDAIV4aRH/VnleeaO8qxsTxTv5a2Yt4h2hSEHhOBdA30y8xjWY6/QUsG
2bX/Fc77gxED23cN54x1eB+otRSqOBOlwS03yCSXd+uC9pLiVUjsCakDeHA4FWA+
uLaVVYsXgT0yPuzVVE61dn1mYMr18jtnsp8XAoqkJlzxK67HO2MvKxlZcy3hbsFj
ikeObha1VeWeUMopDMujPFujsQ+HDyFJt3qdip5a4EpDMSqWzto4tJ+euxTTDhV7
/t33dlhCt7whbjAQhvSzJUM7d5UzA6vFZ9SIfx/gZV9D7Pk4vqdOuy6gTblBg+DQ
UHZL3URngyJXeOznladruf8dAUSRws2Glml/+SvZTQ73q+cB1NEPgkj8CkVNEKcM
gX0hdB7Bdfdwz/tW8X1XimZZgfjtYBG3IWLRLNNMFAh+1wLEfisY9iwabjMbeuXL
Hk/4U6StMwyhu4sSZYnJNnu47t0uZKmul1gyECnnXyvZJVDhO0Ibsq4YWQxz7ogX
oV0lG4sBquYLFewb4YVhH5XywoPgF6g64uj/2NA4VHSXE58BapqQw1JTIQARAQAB
tD50ZmxpbnQtcnVsZXNldC10cmFpbGluZy1jb21tYSByZWxlYXNlIDxyZWxlYXNl
QG5pdGl2ZS5pbnZhbGlkPokCVwQTAQoAQRYhBLds/D1wvNZc+IVc/zAKWqcmMpXV
BQJpxRWuAhsDBQkDwmcABQsJCAcCAiICBhUKCQgLAgQWAgMBAh4HAheAAAoJEDAK
WqcmMpXVy58QAKenTnWEQdBKm2bPsLeuNBzpxDSDdPLECTfLAQOq7qEsKtTKw6sx
PlBOYjcoaewWG/W8ioXUmAXGzqyJMJcZgAcmmaDg2+LeAebDeq2Qpk7E3pNk6+p7
Whe2w9uFhuhmwH5m5laIN6jmNUnsX+i/fNVWY3Y4avJgV2PJqRN4yz4IqNgsEkUt
t26p51gxdp+Nv153iKcK/FCMpS6vvitCicemVkJcCNEwyS7NhQpJR0M0TX8mfhZv
S30zfqkl7R1NAOmUeM0+14WBwuScgAjC8uryZn7dx7jbpnW1RAjKqDHXl80Oh2/5
Nx6cN5KSggYrv/L08+WEKQMuzKwyYL4ADKM/CeHgoMqcXJmpe9MiDNw6Mwk/53dD
jBbJBKesGIc6MSwVwD06j+r/Go/w6mxoSCp6JqgfRiVoPg0FMWGcv2H+ErMoX0IO
XLNvgER9DLh6kpJWgngWZZDp/lC+DrXs4SLpW3z85F2eeXK0aGS2Hat5vqfGzdFt
7nruoYtp5O8S1vu6+Xree1svk/EQ75Fl9qf9+cs4t8ys8sFFGJw9K9jREyt0bvPr
Lb3Dhb7WJ0EbF0LJnsDNaCDkLuZfFATcvKPj+thakcyYUs72qj0f5WkA/Znm+GAD
hKutnByLnk/WXj30RwjnPAzvfMnumsxRC34JiH7olGmQ1wY/CbhgcuRG
=tFFX
-----END PGP PUBLIC KEY BLOCK-----
KEY
}
```

`signing_key` is used when TFLint falls back to PGP signature verification. For private GitHub repositories, set it explicitly.

## Rules

| Name                     | Description                                                                                                | Severity | Enabled | Link |
| ------------------------ | ---------------------------------------------------------------------------------------------------------- | -------- | ------- | ---- |
| multiline_trailing_comma | Enforce trailing commas for multiline lists and function calls, while forbidding them on single-line forms | ERROR    | ✔      |      |
| multiline_map_no_comma   | Forbid commas in multiline maps and trailing commas in single-line maps                                    | ERROR    | ✔      |      |

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
