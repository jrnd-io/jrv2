// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package emitter

import (
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "describes available emitters",
	Long: "describes available emitters. Example usage:\n" +
		"jr emitter list",
	Run: list,
}

func list(cmd *cobra.Command, _ []string) {
	noColor, _ := cmd.Flags().GetBool("nocolor")
	all, _ := cmd.Flags().GetBool("full")

	var Green = ""
	var Reset = ""
	if !noColor {
		Green = "\033[32m"
		Reset = "\033[0m"
	}

	fmt.Println()
	fmt.Println("List of JR emitters:")
	fmt.Println()

	for k, v := range config.Emitters {
		if all {
			fmt.Printf("%s%s%s", Green, k, Reset)
			fmt.Print(" -> (")
			for _, e := range v {
				fmt.Printf(" %s ", e.Name)
			}
			fmt.Println(")")

		} else {
			fmt.Printf("%s%s%s\n", Green, k, Reset)
		}
	}
	fmt.Println()
}

func init() {
	ListCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
	ListCmd.Flags().BoolP("full", "f", false, "Show list of templates for each emitter")
}
