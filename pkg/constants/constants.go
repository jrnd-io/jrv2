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

package constants

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
)

var JRSystemDir string
var JRUserDir string
var SystemDir = fmt.Sprintf("%s%c%s", xdg.ConfigHome, os.PathSeparator, "jr")
var UserDir = fmt.Sprintf("%s%c%s", xdg.DataHome, os.PathSeparator, "jr")

const (
	Num                       = 1
	Locale                    = "us"
	Frequency                 = -1
	Infinite                  = 1<<63 - 1
	DefaultKey                = "null"
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
	DefaultTopic              = "test"
	DefaultHTTPPort           = 7482 //JR :)
	DefaultLogLevel           = "fatal"
)
