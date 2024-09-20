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

package function

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/rs/zerolog/log"
)

var defaultLocale = "us"
var data = map[string][]string{}

// ClearCache is used to internally Cache data from word files
func ClearCache(name string) {
	data[name] = nil
}
func GetCache(name string) []string {
	return data[name]
}
func Cache(name string) (bool, error) {

	v := data[name]
	if v != nil {
		return false, nil
	}
	templateDir := fmt.Sprintf("%s/%s", config.JrSystemDir, "templates")

	locale := emitter.GetState().Locale
	fileName := fmt.Sprintf("%s%cdata%c%s%c%s",
		os.ExpandEnv(templateDir),
		os.PathSeparator,
		os.PathSeparator,
		locale,
		os.PathSeparator,
		name)
	if locale != defaultLocale && !(fileExists(fileName)) {
		fileName = fmt.Sprintf("%s%cdata%c%s%c%s",
			os.ExpandEnv(templateDir),
			os.PathSeparator,
			os.PathSeparator,
			defaultLocale,
			os.PathSeparator,
			name)
	}

	return CacheFromFile(name, fileName)

}

func CacheFromFile(fileName string, name string) (bool, error) {

	data[name] = initialize(fileName)
	if len(data[name]) == 0 {
		return false, fmt.Errorf("no words found in %s", fileName)
	}

	return true, nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	return false
}

func initialize(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error in closing file")
		}
	}(file)

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words
}

func InitCSV(csvpath string) {
	// Loads the csv file in the context
	if len(csvpath) == 0 {
		return
	}

	var csvHeaders = make(map[int]string)
	csvValues := make(emitter.CSVMap)

	if _, err := os.Stat(csvpath); err != nil {
		println("File does not exist: ", csvpath)
		os.Exit(1)
	}

	file, err := os.Open(csvpath)

	if err != nil {
		println("Error opening file:", csvpath, "error:", err)
		os.Exit(1)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		file.Close()
		log.Fatal().Err(err).Str("file", csvpath).Msg("Error reading CSV file") //nolint do not bother
	}

	for row := 0; row < len(lines); row++ {
		var aRow = lines[row]
		for col := 0; col < len(aRow); col++ {
			var value = aRow[col]
			// print(" ROW -> ",row)
			if row == 0 {
				// print(" H: ", value)
				csvHeaders[col] = strings.Trim(value, " ")
			} else {

				val, exists := csvValues[row-1]
				if exists {
					val[csvHeaders[col]] = strings.Trim(value, " ")
					csvValues[row-1] = val
				} else {
					var localmap = make(map[string]string)
					localmap[csvHeaders[col]] = strings.Trim(value, " ")
					csvValues[row-1] = localmap
				}
				//	print(" V: ", value)
			}
		}
		// println()
	}

	emitter.GetState().SetCSV(csvValues)
}
