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
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/jrnd-io/jrv2/pkg/api"
	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show [template]",
	Short: "Show a template",
	Long:  `Show a template. Templates must be in system or in user directory, which are '$JR_USER_DIR/templates' and '$JR_SYSTEM_DIR/templates'`,
	Args:  cobra.ExactArgs(1),
	RunE:  show,
}

func show(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("missing template name: try the list command to see available templates")
	}

	cyan := color.New(color.FgCyan)
	noColor, _ := cmd.Flags().GetBool("nocolor")

	if noColor {
		cyan.DisableColor()
	}
	cyanf := cyan.Sprintf
	templateString, err := api.GetRawValidatedTemplate(args[0])

	if runtime.GOOS != "windows" {
		templateString = strings.ReplaceAll(templateString, "{{", cyanf("{{"))
		templateString = strings.ReplaceAll(templateString, "}}", cyanf("}}"))
	}
	fmt.Println(templateString)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	ShowCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
}
