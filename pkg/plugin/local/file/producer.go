// Copyright © 2024 JR team
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

package console

import (
	"context"
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	PluginName = "file"
	OutputDir  = "outputDir"
	FileName   = "fileName"
)

func init() {
	plugin.RegisterLocalPlugin(PluginName, &plugin.Plugin{
		Name:     PluginName,
		Producer: &Producer{},
		IsRemote: false,
	})
}

type Producer struct {
}

func (p *Producer) Produce(_ context.Context, _ []byte, v []byte, _ map[string]string, configParams map[string]string) (*jrpc.ProduceResponse, error) {

	outputDir := "."
	if configParams[OutputDir] != "" {
		outputDir = configParams[OutputDir]
	}

	fileName := configParams["emitter.name"]
	if configParams[FileName] != "" {
		fileName = configParams[FileName]
	}

	filePath := fmt.Sprintf("%s%c%s", outputDir, os.PathSeparator, fileName)
	log.Debug().Str("file", filePath).Msg("writing to file")
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	nWritten, err := f.Write(v)
	if err != nil {
		log.Warn().Str("file", fileName).Err(err).Msg("failed to write to file")
	}

	return &jrpc.ProduceResponse{
		Bytes:   uint64(nWritten), //nolint
		Message: "",
	}, nil

}

func (p *Producer) Close(_ context.Context) error {
	return nil
}
