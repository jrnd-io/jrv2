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
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strings"
)

var ShowCmd = &cobra.Command{
	Use:   "show [template]",
	Short: "Show a template",
	Long:  `Show a template. Templates must be in system or in user directory, which are '$JR_USER_DIR/templates' and '$JR_SYSTEM_DIR/templates'`,
	Run:   show,
}

func show(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Error().Msg("Missing template name: try the list command to see available templates")
		os.Exit(1)
	}

	noColor, _ := cmd.Flags().GetBool("nocolor")
	templateString, validError := api.GetRawValidatedTemplate(args[0])

	var Reset = "\033[0m"
	if runtime.GOOS != "windows" && !noColor {
		var Cyan = "\033[36m"
		coloredOpeningBracket := fmt.Sprintf("%s%s", Cyan, "{{")
		coloredClosingBracket := fmt.Sprintf("%s%s", "}}", Reset)
		templateString = strings.ReplaceAll(templateString, "{{", coloredOpeningBracket)
		templateString = strings.ReplaceAll(templateString, "}}", coloredClosingBracket)
	}
	fmt.Println(templateString)
	fmt.Print(Reset)
	if validError != nil {
		log.Fatal().Msg("Invalid template")
	}
}

func init() {
	ShowCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
}
