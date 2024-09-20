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
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/adrg/xdg"
)

var JrSystemDir string
var JrUserDir string
var JrSeed int64
var Random *rand.Rand

var SystemDir = fmt.Sprintf("%s%c%s", xdg.DataDirs[0], os.PathSeparator, "jr")
var UserDir = fmt.Sprintf("%s%c%s", xdg.DataHome, os.PathSeparator, "jr")

const (
	DefaultNum                = 1
	DefaultLocale             = "us"
	DefaultFrequency          = "1s"
	DefaultDuration           = "2190000h" // 250 years
	DefaultKeyTemplate        = "null"
	DefaultOutput             = "stdout"
	DefaultOutputTemplate     = "{{.V}}\n"
	DefaultOutputKcatTemplate = "{{.K}},{{.V}}\n"
	KafkaConfig               = "./kafka/config.properties"
	DefaultPartitions         = 6
	DefaultReplica            = 3
	DefaultPreloadSize        = 0
	DefaultEnvPrefix          = "JR"
	DefaultEmitterName        = "cli"
	DefaultValueTemplate      = "user"
	DefaultHeaderTemplate     = "null"
	DefaultTopic              = "test"
	DefaultHTTPPort           = 7482 // JR :)
	DefaultLogLevel           = "fatal"
)

func init() {
	initEnvironmentVariables()
}

func initEnvironmentVariables() {
	JrSystemDir = os.Getenv("JR_SYSTEM_DIR")
	JrUserDir = os.Getenv("JR_USER_DIR")
	seed, err := strconv.ParseInt(os.Getenv("JR_SEED"), 10, 64)

	if (seed == -1) || (err != nil) {
		JrSeed = time.Now().UTC().UnixNano()
	} else {
		JrSeed = seed
	}

	Random = rand.New(rand.NewSource(JrSeed)) //nolint no need for a secure random generator

	if JrSystemDir == "" {
		JrSystemDir = SystemDir
	}
	if JrUserDir == "" {
		JrUserDir = UserDir
	}
}
