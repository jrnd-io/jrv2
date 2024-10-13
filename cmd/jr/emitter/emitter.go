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
	"os"
	"strings"

	"github.com/jrnd-io/jrv2/pkg/config"
	emitterapi "github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var emitterCmd = &cobra.Command{
	Use:     "emitter",
	Short:   "jr Emitter resource",
	Long:    `jr Emitter resource`,
	GroupID: "resource",
}

func NewCmd() *cobra.Command {
	emitterCmd.AddCommand(ListCmd)
	emitterCmd.AddCommand(RunCmd)
	emitterCmd.AddCommand(ShowCmd)
	return emitterCmd
}

func initEmitters() {
	viper.SetConfigName("jrconfig")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	for _, path := range strings.Split(os.ExpandEnv("$PATH"), ":") {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath(config.JrSystemDir)

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("JR configuration not found")
	}
	err := viper.UnmarshalKey("Emitters", &emitterapi.Emitters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal emitter configuration")
	}
}

func init() {
	config.InitEnvironmentVariables()
	initEmitters()
}
