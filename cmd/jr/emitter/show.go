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

	"github.com/fatih/color"
	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show a configured emitter",
	Long: `show a configured emitter:
jr emitter show net_device`,
	Args: cobra.ExactArgs(1),
	Run:  show,
}

func show(cmd *cobra.Command, args []string) {

	noColor, _ := cmd.Flags().GetBool("nocolor")

	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)
	if noColor {
		green.DisableColor()
		white.DisableColor()
	}

	greenf := green.Printf
	whitef := white.Sprintf

	fmt.Println()
	for k, v := range emitter.Emitters {
		if k == args[0] {
			for _, e := range v {
				greenf("Name: %s\n", whitef("%s", e.Name))                //nolint
				greenf("Locale: %s\n", whitef("%s", e.Locale))            //nolint
				greenf("Num: %s\n", whitef("%d", e.Tick.Num))             //nolint
				greenf("Frequency: %s\n", whitef("%d", e.Tick.Frequency)) //nolint
				greenf("Duration: %s\n", whitef("%d", e.Tick.Duration))   //nolint
				greenf("Preload: %s\n", whitef("%d", e.Preload))          //nolint
				greenf("Output: %s\n", whitef(e.Output))                  //nolint
				greenf("Oneline: %s\n", whitef("%b", e.Oneline))          //nolint
				greenf("Key Template: %s\n", whitef(e.KeyTemplate))       //nolint
				greenf("Value Template: %s\n", whitef(e.ValueTemplate))   //nolint
				greenf("Output Template: %s\n", whitef(e.OutputTemplate)) //nolint
			}
		}
	}
	fmt.Println()

}

func init() {
	ShowCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
}
