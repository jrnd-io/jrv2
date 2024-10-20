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

	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/tpl"

	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template.
  Without any other flag, [template] is just the name of a template in the templates directories, which are '$JR_SYSTEM_DIR/templates' and '$JR_USER_DIR/templates'. Example:
jr template run net_device
  With the --embedded flag, [template] is a string containing a full template. Example:
jr template run --embedded "{{name}}"
`,
	Args: cobra.ExactArgs(1),
	RunE: run,
}

func run(cmd *cobra.Command, args []string) error {
	keyTemplate, _ := cmd.Flags().GetString("key")
	headerTemplate, _ := cmd.Flags().GetString("header")

	outputTemplate, _ := cmd.Flags().GetString("outputTemplate")
	embedded, _ := cmd.Flags().GetBool("embedded")
	kcat, _ := cmd.Flags().GetBool("kcat")
	output, _ := cmd.Flags().GetString("output")
	oneline, _ := cmd.Flags().GetBool("oneline")
	locale, _ := cmd.Flags().GetString("locale")

	num, _ := cmd.Flags().GetInt("num")
	frequency, _ := cmd.Flags().GetDuration("frequency")
	duration, _ := cmd.Flags().GetDuration("duration")
	immediateStart, _ := cmd.Flags().GetBool("immediate")
	throughputString, _ := cmd.Flags().GetString("throughput")
	preload, _ := cmd.Flags().GetInt("preload")

	// csv, _ := cmd.Flags().GetString("csv")

	if kcat {
		oneline = true
		output = "stdout"
		outputTemplate = config.DefaultOutputKcatTemplate
	}

	var valueTemplate string
	var err error
	if embedded {
		valueTemplate = args[0]
	} else {
		valueTemplate, err = tpl.GetRawTemplate(args[0])
		if err != nil {
			return err
		}
	}

	if throughputString != "" {
		result, err := tpl.ExecuteTemplate(valueTemplate, nil)
		if err != nil {
			return err
		}
		throughput, err := emitter.ParseThroughput(throughputString)
		if err != nil {
			return err
		}
		frequency = emitter.CalculateFrequency(len([]byte(result)), num, throughput)
	}

	e, err := emitter.New(
		emitter.WithDuration(duration),
		emitter.WithFrequency(frequency),
		emitter.WithNum(num),
		emitter.WithPreload(preload),
		emitter.WithKeyTemplate(keyTemplate),
		emitter.WithValueTemplate(valueTemplate),
		emitter.WithHeaderTemplate(headerTemplate),
		emitter.WithOutputTemplate(outputTemplate),
		emitter.WithImmediateStart(immediateStart),
		emitter.WithOutput(output),
		emitter.WithLocale(locale),
		emitter.WithOneline(oneline),
	)
	if err != nil {
		return err
	}
	fmt.Println(e)
	return nil
}

func init() {
	RunCmd.Flags().IntP("num", "n", config.DefaultNum, "Number of elements to create for each pass")
	frequency := config.DefaultFrequency
	RunCmd.Flags().DurationP("frequency", "f", frequency, "how much time to wait for next generation pass")
	duration := config.DefaultDuration
	RunCmd.Flags().DurationP("duration", "d", duration, "If frequency is enabled, with Duration you can set a finite amount of time")
	RunCmd.Flags().String("throughput", "", "You can set throughput, JR will calculate frequency automatically.")
	RunCmd.Flags().Int("preload", config.DefaultPreloadSize, "Number of elements to create during the preload phase")
	RunCmd.Flags().Bool("embedded", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")
	RunCmd.Flags().Bool("immediate", false, "If frequency is enabled, it will tick immediately too")
	RunCmd.Flags().StringP("key", "k", config.DefaultKeyTemplate, "A template to generate a key")
	RunCmd.Flags().String("header", config.DefaultHeaderTemplate, "A template to generate a header")
	RunCmd.Flags().StringP("output", "o", config.DefaultOutput, "can be one of stdout, kafka, http, redis, mongo, elastic, s3, gcs, azblobstorage, azcosmosdb, cassandra, luascript, wasm, awsdynamodb")
	RunCmd.Flags().String("outputTemplate", config.DefaultOutputTemplate, "Formatting of K,V on standard output")
	RunCmd.Flags().BoolP("oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
	RunCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
	RunCmd.Flags().String("locale", config.DefaultLocale, "DefaultLocale")
	RunCmd.Flags().String("csv", "", "Path to csv file to use")
}
