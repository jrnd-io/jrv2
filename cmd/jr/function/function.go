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
	"github.com/spf13/cobra"
)

var functionCmd = &cobra.Command{
	Use:     "function",
	Short:   "jr Function resource",
	Long:    `jr Function resource`,
	GroupID: "resource",
}

func NewListCmd() *cobra.Command {
	ListCmd.Flags().BoolP("category", "c", false, "IndexOf in category")
	ListCmd.Flags().BoolP("find", "f", false, "IndexOf in description and name")
	ListCmd.Flags().BoolP("markdown", "m", false, "Output the list as markdown")
	ListCmd.Flags().BoolP("run", "r", false, "Run the example")
	ListCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
	functionCmd.AddCommand(ListCmd)
	return functionCmd
}

func NewManCmd() *cobra.Command {
	ManCmd.Flags().BoolP("category", "c", false, "IndexOf in category")
	ManCmd.Flags().BoolP("find", "f", false, "IndexOf in description and name")
	ManCmd.Flags().BoolP("markdown", "m", false, "Output the list as markdown")
	ManCmd.Flags().BoolP("run", "r", false, "Run the example")
	ManCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
	return ManCmd
}
