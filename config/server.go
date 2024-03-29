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

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/ccx-data-pipeline-monitor/packages/config/server.html

import (
	"github.com/spf13/viper"
)

// ServerConfig data type represents configuration of server
type ServerConfig struct {
	Address  string
	UseHTTPS bool
}

// ReadServerConfig function reads configuration options for HTTP server with this service
func ReadServerConfig() ServerConfig {
	var cfg ServerConfig
	sub := viper.Sub("server")
	cfg.Address = sub.GetString("address")
	cfg.UseHTTPS = sub.GetBool("use_https")
	return cfg
}
