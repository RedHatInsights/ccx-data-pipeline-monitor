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

package config

import (
	"github.com/spf13/viper"
)

// OpenShiftConfig represents all configuration options required to get access to OpenShift via oc client
type OpenShiftConfig struct {
	URL     string
	Project string
}

// ReadOpenShiftConfig function reads configuration options required to get access to OpenShift via oc client
func ReadOpenShiftConfig() OpenShiftConfig {
	var cfg OpenShiftConfig
	sub := viper.Sub("openshift")
	cfg.URL = sub.GetString("url")
	cfg.Project = sub.GetString("project")
	return cfg
}
