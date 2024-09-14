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
	"github.com/jrnd-io/jrv2/pkg/constants"
	"github.com/spf13/cobra"
	"time"
)

var RunCmd = &cobra.Command{
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template.
  Without any other flag, [template] is just the name of a template in the templates directory, which is '$JR_SYSTEM_DIR/templates'. Example:
jr template run net_device
  With the --embedded flag, [template] is a string containing a full template. Example:
jr template run --template "{{name}}"
`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {

}

func init() {
	RunCmd.Flags().IntP("num", "n", constants.Num, "Number of elements to create for each pass")
	RunCmd.Flags().DurationP("frequency", "f", constants.Frequency, "how much time to wait for next generation pass")
	RunCmd.Flags().DurationP("duration", "d", constants.Infinite, "If frequency is enabled, with Duration you can set a finite amount of time")
	RunCmd.Flags().String("throughput", "", "You can set throughput, JR will calculate frequency automatically.")
	RunCmd.Flags().Int64("seed", time.Now().UTC().UnixNano(), "Seed to init pseudorandom generator")
	RunCmd.Flags().String("csv", "", "Path to csv file to use")
	RunCmd.Flags().Bool("embedded", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")
	RunCmd.Flags().StringP("key", "k", constants.DefaultKey, "A template to generate a key")

	/*
		templateRunCmd.Flags().StringP("kafkaConfig", "F", "", "Kafka configuration")
		templateRunCmd.Flags().String("registryConfig", "", "Kafka configuration")
		templateRunCmd.Flags().Int("preload", constants.DEFAULT_PRELOAD_SIZE, "Number of elements to create during the preload phase")

		templateRunCmd.Flags().StringP("topic", "t", constants.DEFAULT_TOPIC, "Kafka topic")

		templateRunCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
		templateRunCmd.Flags().StringP("output", "o", constants.DEFAULT_OUTPUT, "can be one of stdout, kafka, http, redis, mongo, elastic, s3, gcs, azblobstorage, azcosmosdb, cassandra, luascript, wasm, awsdynamodb")
		templateRunCmd.Flags().String("outputTemplate", constants.DEFAULT_OUTPUT_TEMPLATE, "Formatting of K,V on standard output")
		templateRunCmd.Flags().BoolP("oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
		templateRunCmd.Flags().BoolP("autocreate", "a", false, "if enabled, autocreate topics")
		templateRunCmd.Flags().String("locale", constants.LOCALE, "Locale")

		templateRunCmd.Flags().BoolP("schemaRegistry", "s", false, "If you want to use Confluent Schema Registry")
		templateRunCmd.Flags().String("serializer", "", "Type of serializer: json-schema, avro-generic, avro, protobuf")
		templateRunCmd.Flags().Duration("redis.ttl", -1, "If output is redis, ttl of the object")
		templateRunCmd.Flags().String("httpConfig", "", "HTTP configuration")
		templateRunCmd.Flags().String("redisConfig", "", "Redis configuration")
		templateRunCmd.Flags().String("mongoConfig", "", "MongoDB configuration")
		templateRunCmd.Flags().String("elasticConfig", "", "Elastic Search configuration")
		templateRunCmd.Flags().String("s3Config", "", "AWS S3 configuration")
		templateRunCmd.Flags().String("awsDynamoDBConfig", "", "AWS DynamoDB configuration")
		templateRunCmd.Flags().String("gcsConfig", "", "Google GCS configuration")
		templateRunCmd.Flags().String("azBlobStorageConfig", "", "Azure Blob storage configuration")
		templateRunCmd.Flags().String("azCosmosDBConfig", "", "Azure CosmosDB configuration")
		templateRunCmd.Flags().String("cassandraConfig", "", "Cassandra configuration")
		templateRunCmd.Flags().String("luascriptConfig", "", "LUA Script configuration")
		templateRunCmd.Flags().String("wasmConfig", "", "WASM configuration")
	*/
}
