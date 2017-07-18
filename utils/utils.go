// Copyright © 2017 Nico Balestra <functions@protonmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/spf13/viper"
)

func tabularize(elems []string) string {
	toReturn := ""

	for _, val := range elems {
		toReturn = toReturn + val + "\t"
	}

	return toReturn
}

func PrintTabular(headers []string, data [][]string) {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.FilterHTML)
	defer w.Flush()
	fmt.Println("")
	headersStr := tabularize(headers)
	fmt.Fprintln(w, headersStr)

	for _, row := range data {
		lineStr := tabularize(row)
		fmt.Fprintln(w, lineStr)
	}

}

func Debug(doSomething func()) {
	if viper.GetBool(KEY_DEBUG) {
		doSomething()
	}
}
