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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()
	fmt.Println("")
	headersStr := tabularize(headers)
	fmt.Fprintln(w, headersStr)

	for _, row := range data {
		lineStr := tabularize(row)
		fmt.Fprintln(w, lineStr)
	}

}

type DebugFunc func()

func Debug(doSomething DebugFunc) {
	if viper.GetBool(KEY_DEBUG) {
		doSomething()
	}
}
