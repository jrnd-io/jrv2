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
)

const (
	PluginName = "stdout"
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

func (p *Producer) Produce(_ context.Context, _ []byte, v []byte, _ map[string]string, _ map[string]string) (*jrpc.ProduceResponse, error) {

	writtenBytes := len(v)
	fmt.Println(string(v))
	return &jrpc.ProduceResponse{
		Bytes:   uint64(writtenBytes),
		Message: "",
	}, nil

}

func (p *Producer) Close(_ context.Context) error {
	return nil
}
