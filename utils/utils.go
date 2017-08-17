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
	"strings"
	"encoding/json"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const (
	tabwriterMinWidth = 10
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = 0
)

const (
	OUTPUT_JSON = "json"
	OUTPUT_YAML = "yaml"
	OUTPUT_LIST = "list"
)

func tabularize(elems []string) string {
	toReturn := ""

	for _, val := range elems {
		toReturn = toReturn + strings.TrimSpace(val) + "\t"
	}

	return toReturn
}

func PrintTabular(headers []string, data [][]string) {

	w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
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

func PrintObject(objects []interface{}, headers []string, extractTabularDataFunc func([]interface{}) [][]string) {

	switch viper.GetString(KEY_FORMAT) {
	case OUTPUT_JSON:
		PrintAsJSON(objects)
		break;
	case OUTPUT_YAML:
		panic("YAML output not implemented yet")
	case OUTPUT_LIST:
		PrintAsList(objects, extractTabularDataFunc, headers)
	default:
		PrintAsList(objects, extractTabularDataFunc, headers)
	}
}

func PrintAsList(objects []interface{}, extractTabularDataFunc func([]interface{}) [][]string, headers []string) {

	data := extractTabularDataFunc(objects)

	PrintTabular(headers, data)
}

func PrintAsJSON(objects []interface{}) {
	b, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		fmt.Println("Error while marshalling output:", err)
	}
	os.Stdout.Write(b)
}

func OpenYAMLFile(f string, t interface{}) (error) {

	fileContent, err := ioutil.ReadFile(f)

	if err != nil {
		return fmt.Errorf("Error while opening yaml file %q. Error: %s", f, err)
	}

	err = yaml.Unmarshal(fileContent, t)

	if err != nil {
		return fmt.Errorf("Error while parsing YAML file %q . Error: %s", f, err)
	}

	return nil
}



