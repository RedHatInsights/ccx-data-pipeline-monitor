# ccx-data-pipeline-monitor

[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)

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
* [BDD tests](#bdd-tests)
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



## BDD tests

Behaviour tests for this service are included in [Insights Behavioral
Spec](https://github.com/RedHatInsights/insights-behavioral-spec) repository.
In order to run these tests, the following steps need to be made:

1. clone the [Insights Behavioral Spec](https://github.com/RedHatInsights/insights-behavioral-spec) repository
1. go into the cloned subdirectory `insights-behavioral-spec`
1. run the `ccx_data_pipeline_monitor_tests.sh` from this subdirectory

List of all test scenarios prepared for this service is available at
<https://github.com/RedHatInsights/insights-behavioral-spec#ccx-data-pipeline-monitor>



## Package manifest

Package manifest is available at [docs/manifest.txt](docs/manifest.txt).
