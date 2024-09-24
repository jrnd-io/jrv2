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
	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/random"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "prints JR environment variables",
	Long:  `prints JR environment variables`,
	Run: func(_ *cobra.Command, _ []string) {

		seed := "AUTO"
		randomGen := "Global default (not reproducible, thread safe)"
		if random.JrSeed != -1 {
			seed = fmt.Sprintf("%x", random.CreateByteSeed(uint64(random.JrSeed)))
			randomGen = "Local ChaCha8 (reproducible, not thread safe)"
		}
		fmt.Printf("JR System Dir: %s\n", config.JrSystemDir)
		fmt.Printf("JR User Dir  : %s\n", config.JrUserDir)
		fmt.Printf("JR Seed      : %s\n", seed)
		fmt.Printf("JR conf file : %s\n", viper.ConfigFileUsed())
		fmt.Printf("Random Gen   : %s\n", randomGen)
		fmt.Printf("Log Level    : %s\n", config.LogLevel)
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
