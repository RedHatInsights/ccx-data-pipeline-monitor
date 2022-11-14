# ccx-data-pipeline-monitor
[![GoDoc](https://godoc.org/github.com/RedHatInsights/ccx-data-pipeline-monitor?status.svg)](https://godoc.org/github.com/RedHatInsights/ccx-data-pipeline-monitor)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://redhatinsights.github.io/ccx-data-pipeline-monitor/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/ccx-data-pipeline-monitor)](https://goreportcard.com/report/github.com/RedHatInsights/ccx-data-pipeline-monitor)
[![Build Status](https://travis-ci.org/RedHatInsights/ccx-data-pipeline-monitor.svg?branch=master)](https://travis-ci.org/RedHatInsights/ccx-data-pipeline-monitor)
[![codecov](https://codecov.io/gh/RedHatInsights/ccx-data-pipeline-monitor/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/ccx-data-pipeline-monitor)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/RedHatInsights/ccx-data-pipeline-monitor)
[![License](https://img.shields.io/badge/license-Apache-blue)](https://github.com/RedHatInsights/ccx-data-pipeline-monitor/blob/master/LICENSE)

Monitor for CCX data pipeline

<!-- vim-markdown-toc GFM -->

* [Makefile targets](#makefile-targets)
* [Package manifest](#package-manifest)

<!-- vim-markdown-toc -->

## Makefile targets

```
Usage: make <OPTIONS> ... <TARGETS>

Available targets are:

clean                Run go clean
build                Run go build
fmt                  Run go fmt -w for all sources
lint                 Run golint
vet                  Run go vet. Report likely mistakes in source code
cyclo                Run gocyclo
ineffassign          Run ineffassign checker
shellcheck           Run shellcheck
errcheck             Run errcheck
goconst              Run goconst checker
style                Run all the formatting related commands (fmt, vet, lint, cyclo) + check shell scripts
run                  Build the project and executes the binary
test                 Run the unit tests
help                 Show this help screen
function_list        List all functions in generated binary file
```

## Package manifest

Package manifest is available at [docs/manifest.txt](docs/manifest.txt).
