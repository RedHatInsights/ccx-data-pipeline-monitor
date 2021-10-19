/*
Copyright © 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oc

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/ccx-data-pipeline-monitor/packages/oc/oc.html

import (
	"bytes"
	"os/exec"
	"strings"
)

// Command run any oc command and return its standard and error outputs
func Command(args ...string) (string, string, error) {
	// disable "G204 (CWE-78): Subprocess launched with variable
	// #nosec G204
	cmd := exec.Command("oc", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	outString := string(stdout.Bytes())
	errString := string(stderr.Bytes())
	if err != nil {
		return outString, errString, err
	}

	return outString, errString, nil
}

// Login perform login into oc
func Login(url, arg string) (string, string, error) {
	token := getToken(arg)
	return Command("login", url, "--token="+token)
}

// GetPods function reads list of all pods via oc command
func GetPods() (string, string, error) {
	return Command("get", "pods")
}

// GetLogs functions reads logs for selected pod
func GetLogs(pod string) (string, string, error) {
	return Command("logs", pod)
}

func getToken(arg string) string {
	const tokenPart = "--token="

	token := arg

	// check whether just token is provided or the whole oc login command
	i := strings.LastIndex(arg, tokenPart)
	if i >= 0 && len(arg) >= i+len(tokenPart) {
		// get just the token part
		token = arg[i+len(tokenPart):]
	}

	return token
}
