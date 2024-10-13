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

package main

import (
	"fmt"
	"os"

	"github.com/jrnd-io/jrv2/cmd/jr/emitter"
	"github.com/jrnd-io/jrv2/cmd/jr/function"
	"github.com/jrnd-io/jrv2/cmd/jr/producer"
	"github.com/jrnd-io/jrv2/cmd/jr/template"
	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jr",
	Short: "jr, the data random generator",
	Long:  `jr is a data random generator that helps you in generating quality random data for your needs.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddGroup(&cobra.Group{
		ID:    "resource",
		Title: "Resources",
	})

	rootCmd.PersistentFlags().StringVar(&config.JrSystemDir, "jr_system_dir", config.JrSystemDir, "JR system dir")
	rootCmd.PersistentFlags().StringVar(&config.JrUserDir, "jr_user_dir", config.JrUserDir, "JR user dir")
	rootCmd.PersistentFlags().CountP("verbose", "v", "Set the verbosity level")

	rootCmd.AddCommand(emitter.NewCmd())
	rootCmd.AddCommand(producer.NewCmd())
	rootCmd.AddCommand(function.NewListCmd())
	rootCmd.AddCommand(function.NewManCmd())
	rootCmd.AddCommand(template.NewCmd())

	config.InitEnvironmentVariables()
}

func initConfig() {
	// Get the verbosity level
	verbosity, err := rootCmd.PersistentFlags().GetCount("verbose")
	if err != nil || verbosity == 0 {
		verbosity = config.JrVerbosity
	}

	// setting zerolog level
	var zlogLevel zerolog.Level
	switch verbosity {
	case 1:
		zlogLevel = zerolog.ErrorLevel
	case 2:
		zlogLevel = zerolog.WarnLevel
	case 3:
		zlogLevel = zerolog.InfoLevel
	case 4:
		zlogLevel = zerolog.DebugLevel
	case 5:
		zlogLevel = zerolog.TraceLevel
	default:
		zlogLevel = zerolog.Disabled
	}

	zerolog.SetGlobalLevel(zlogLevel)

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.999Z07:00"}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
