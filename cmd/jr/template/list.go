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
package template

import (
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/api"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Long:  `List all available templates, which are in '$JR_SYSTEM_DIR/templates' and '$JR_USER_DIR/templates' directory`,
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {

	noColor, _ := cmd.Flags().GetBool("nocolor")
	fullPath, _ := cmd.Flags().GetBool("fullPath")

	fmt.Println()
	fmt.Println("System JR templates:")
	fmt.Println()
	printTemplateList(api.SystemTemplateList(), noColor, fullPath)
	fmt.Println()
	fmt.Println("User JR templates:")
	fmt.Println()
	//printTemplateList(api.UserTemplateList(), noColor, fullPath)

}

func printTemplateList(templateList []*api.TemplateInfo, noColor bool, fullPath bool) {

	if len(templateList) == 0 {
		return
	}

	var Red = "\033[31m"
	var Green = "\033[32m"
	var Reset = "\033[0m"

	for _, t := range templateList {
		if !noColor {
			if t.IsValid {
				fmt.Print(Green)
			} else {
				fmt.Print(Red)
			}
		}

		if fullPath {
			fmt.Println(t.FullPath)
		} else {
			fmt.Print(t.Name)
			if t.Error != nil {
				fmt.Println(" -> ", t.Error)
			} else {
				fmt.Println()
			}
		}
	}
	if !(noColor) {
		fmt.Println(Reset)
	}
}
