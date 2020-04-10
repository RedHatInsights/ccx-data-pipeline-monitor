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

package commands

import (
	"fmt"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/oc"
)

func TryToLogin(url string, ocLogin string) bool {
	stdout, stderr, err := oc.Login(url, ocLogin)
	if err == nil {
		fmt.Println(colorizer.Green("\nDone: you have been loged in to OpenShift"))
		return true
	}
	fmt.Println(colorizer.Red("\nUnable to login to OpenShift"))
	fmt.Println(stdout)
	fmt.Println(stderr)
	return false
}
