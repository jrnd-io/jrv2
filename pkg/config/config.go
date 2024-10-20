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

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/adrg/xdg"
	"github.com/jrnd-io/jrv2/pkg/random"
)

var JrSystemDir string
var JrUserDir string
var JrVerbosity int

var DefaultSystemDir = fmt.Sprintf("%s%c%s", xdg.DataDirs[0], os.PathSeparator, "jr")
var DefaultUserDir = fmt.Sprintf("%s%c%s", xdg.DataHome, os.PathSeparator, "jr")

const (
	EnvJRSystemDir = "JR_SYSTEM_DIR"
	EnvJRUserDir   = "JR_USER_DIR"
	EnvJRVerbosity = "JR_VERBOSITY"
	EnvJRSeed      = "JR_SEED"

	DefaultEmitterName        = "cli"
	DefaultLocale             = "us"
	DefaultNum                = 1
	DefaultPreloadSize        = 0
	DefaultFrequency          = time.Second
	DefaultThroughput         = -1
	DefaultDuration           = 2190000 * time.Hour
	DefaultKeyTemplate        = "null"
	DefaultValueTemplate      = "net_device"
	DefaultHeaderTemplate     = "null"
	DefaultOutput             = "stdout"
	DefaultOutputTemplate     = "{{.V}}\n"
	DefaultOutputKcatTemplate = "{{.K}},{{.V}}\n"
	KafkaConfig               = "./kafka/config.properties"
	DefaultPartitions         = 6
	DefaultReplica            = 3
	DefaultTopic              = "test"
	DefaultHTTPPort           = 7482 // JR :)
	DefaultEnvPrefix          = "JR"
	DefaultVerbosity          = 0
)

func InitEnvironmentVariables() {
	JrSystemDir = os.Getenv(EnvJRSystemDir)
	JrUserDir = os.Getenv(EnvJRUserDir)
	JrSeedEnv := os.Getenv(EnvJRUserDir)
	JrVerbosityEnv := os.Getenv(EnvJRVerbosity)
	seed, err := strconv.ParseInt(JrSeedEnv, 10, 64)
	if err != nil || JrSeedEnv == "" {
		seed = -1
	}
	random.SetRandom(seed)

	if JrSystemDir == "" {
		JrSystemDir = DefaultSystemDir
	}
	if JrUserDir == "" {
		JrUserDir = DefaultUserDir
	}
	if JrVerbosityEnv == "" {
		JrVerbosity = DefaultVerbosity
	} else {
		JrVerbosity, err = strconv.Atoi(JrVerbosityEnv)
		if err != nil {
			JrVerbosity = DefaultVerbosity
		}
	}
}
