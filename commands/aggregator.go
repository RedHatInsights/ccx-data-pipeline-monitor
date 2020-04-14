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

	"github.com/c-bata/go-prompt"

	"github.com/RedHatInsights/ccx-data-pipeline-monitor/analyser"
)

func DisplayAggregatorStatistic() {
	fmt.Println(colorizer.Magenta("Aggregator statistic"))
	analyser.PrintAggregatorStatistic(colorizer)
}

func DisplayAggregatorLogs() {
	fmt.Println(colorizer.Magenta("Aggregator logs"))
	fmt.Println(colorizer.Cyan("1."), "consumed but not read")
	fmt.Println(colorizer.Cyan("2."), "read but not whitelisted")
	fmt.Println(colorizer.Cyan("3."), "whitelisted but not marshalled")
	fmt.Println(colorizer.Cyan("4."), "marshalled but not checked")
	fmt.Println(colorizer.Cyan("5."), "checked but not stored")
	fmt.Println()

	which := prompt.Input("selection: ", NoOpCompleter)
	switch which {
	case "1":
		analyser.PrintAggregatorConsumedNotReadMessages(colorizer)
	case "2":
		analyser.PrintAggregatorConsumedNotWhitelisted(colorizer)
	case "3":
		analyser.PrintAggregatorWhitelistedNotMarshalled(colorizer)
	case "4":
		analyser.PrintAggregatorMarshalledNotChecked(colorizer)
	case "5":
		analyser.PrintAggregatorCheckedNotStored(colorizer)
	default:
		fmt.Println(colorizer.Red("wrong input, skipping"))
	}
}
