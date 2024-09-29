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

package producer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	_ "github.com/jrnd-io/jrv2/pkg/plugin/local" //nolint
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "describes available producers",
	Long: "describes available producers. Example usage:\n" +
		"jr producer list",
	Run: list,
}

func list(cmd *cobra.Command, _ []string) {
	noColor, _ := cmd.Flags().GetBool("nocolor")

	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)
	if noColor {
		green.DisableColor()
		white.DisableColor()
	}

	greenf := green.Printf
	whitesf := white.Sprintf

	for _, p := range plugin.GetPluginMap() {
		isRemote := "(local)"
		if p.IsRemote {
			isRemote = "(remote)"
		}

		greenf("%s %s\n", p.Name, whitesf(isRemote)) //nolint
	}

}

func init() {
	config.InitEnvironmentVariables()
	ListCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")

	p, _ := scanForPlugins(fmt.Sprintf("%s%c%s", config.JrSystemDir, os.PathSeparator, "plugins"))
	for _, pl := range p {
		plugin.RegisterRemotePlugin(filepath.Base(pl), &plugin.Plugin{Command: pl})
	}

	p, _ = scanForPlugins(fmt.Sprintf("%s%c%s", config.JrUserDir, os.PathSeparator, "plugins"))
	for _, pl := range p {
		plugin.RegisterRemotePlugin(filepath.Base(pl), &plugin.Plugin{Command: pl})
	}
}

func scanForPlugins(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	plugins := make([]string, 0)
	for _, f := range files {
		if f.Type().IsRegular() && !f.IsDir() {
			plugins = append(plugins, f.Name())
		}
	}

	return plugins, nil

}
