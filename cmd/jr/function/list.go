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

package function

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/spf13/cobra"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "describes available functions",
	Long: "describes available functions. Example usage:\n" +
		"jr function list lorem",
	Run: list,
}

var ManCmd = &cobra.Command{
	Use:   "man",
	Short: "describes available functions",
	Long: "describes available functions. Example usage:\n" +
		"jr function man lorem",
	Run: list,
}

func list(cmd *cobra.Command, args []string) {
	category, _ := cmd.Flags().GetBool("category")
	find, _ := cmd.Flags().GetBool("find")
	run, _ := cmd.Flags().GetBool("run")
	isMarkdown, _ := cmd.Flags().GetBool("markdown")
	noColor, _ := cmd.Flags().GetBool("nocolor")

	switch {
	case category && len(args) > 0:
		funcMap := orderedmap.New[string, function.Description]()
		for k, v := range function.DescriptionMap() {
			if v.Category == args[0] {
				funcMap.Set(k, v)
			}
		}
		printFunctionList(funcMap, isMarkdown, noColor)
	case find && len(args) > 0:
		funcMap := orderedmap.New[string, function.Description]()
		for k, v := range function.DescriptionMap() {
			if strings.Contains(v.Description, args[0]) || strings.Contains(v.Name, args[0]) {
				funcMap.Set(k, v)
			}
		}
		printFunctionList(funcMap, isMarkdown, noColor)
	case len(args) == 1:
		if run {
			f, found := printFunction(args[0], isMarkdown, noColor)
			if found {
				fmt.Println()
				cmd := exec.Command("/bin/sh", "-c", f.Example) // #nosec G204
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			}
		} else {
			printFunction(args[0], isMarkdown, noColor)
		}
	default:
		funcMap := orderedmap.New[string, function.Description]()

		for k, v := range function.DescriptionMap() {
			funcMap.Set(k, v)
		}
		printFunctionList(funcMap, isMarkdown, noColor)

	}

	fmt.Println()
}

func printFunctionList(om *orderedmap.OrderedMap[string, function.Description], isMarkdown bool, noColor bool) {
	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		printFunction(pair.Key, isMarkdown, noColor)
	}
}

/*
func sortAndPrint(functionNames []string, isMarkdown bool, noColor bool) {
	slices.Sort(functionNames)
	for _, k := range functionNames {
		printFunction(k, isMarkdown, noColor)
	}
	fmt.Println()
	fmt.Printf("Total functions: %d\n", len(functionNames))
}
*/

func printFunction(name string, isMarkdown bool, noColor bool) (function.Description, bool) {
	f, found := function.GetDescription(name)

	cyan := color.New(color.FgCyan)
	var white = color.New(color.FgWhite)
	if noColor {
		cyan.DisableColor()
		white.DisableColor()
	}

	cyanf := cyan.Printf
	whitef := white.Sprintf

	if !found {
		return f, found
	}

	if isMarkdown {
		fmt.Println()
		fmt.Printf("### %s \n", f.Name)
		fmt.Printf("**Category:** %s\\\n", f.Category)
		fmt.Printf("**Description:** %s\\\n", f.Description)

		if len(f.Parameters) > 0 {
			fmt.Printf("**Parameters:** `%s`\\\n", f.Parameters)
		} else {
			fmt.Printf("**Parameters:** %s \\\n", f.Parameters)
		}
		fmt.Printf("**Localizable:** `%v`\\\n", f.Localizable)
		fmt.Printf("**Return:** `%s`\\\n", f.Return)
		fmt.Printf("**Example:** `%s`\\\n", f.Example)
		fmt.Printf("**Output:** `%s`\n", f.Output)
	} else {
		fmt.Println()
		cyanf("Name: %s\n", whitef(f.Name)) //nolint
		/*
			cyanf("Category: %s\n", whitef(f.Category))       //nolint
			cyanf("Description: %s\n", whitef(f.Description)) //nolint
			cyanf("Parameters: %s\n", whitef(f.Parameters))   //nolint
			cyanf("Localizable: %s\n", whitef(f.Localizable)) //nolint
			cyanf("Return: %s\n", whitef(f.Return))           //nolint
			cyanf("Example: %s\n", whitef(f.Example))         //nolint
			cyanf("Output: %s\n", whitef(f.Output))           //nolint
		*/
	}

	return f, found

}
