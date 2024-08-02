# Google Chrome Extension ID Generator

[![Go Reference](https://pkg.go.dev/badge/github.com/kadirgun./ceid-generator.svg)](https://pkg.go.dev/github.com/kadirgun./ceid-generator)

This is a simple Go package that generates a Google Chrome Extension ID focuses on finding the first characters.

Note: _Only characters between a-p can be used in Chrome's IDs._

## Installation

```bash
go install github.com/kadirgun/ceid-generator
```

## Usage

```bash
ceid-generator --prefix ceid --threads 4

# Output
# Possibility: 0.001526%
# Estimated tries: 65536
# Tries: 57344
# ID: ceidkbgppdoanomdbbkpojlmicahnlnl
# Public Key: MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2zZ...
```
