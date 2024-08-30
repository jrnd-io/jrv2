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
package wasm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	"github.com/tetratelabs/wazero"
	"os"
	"sync"

	wazapi "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const (
	ScriptName = "wasm"
	PluginName = "wasm"
)

func init() {
	p, err := New()
	if err == nil {
		plugin.RegisterLocalPlugin(PluginName, p)
	}
}

func New() (*Plugin, error) {

	cfgBytes, err := plugin.ReadConfig(PluginName)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := json.Unmarshal(cfgBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal WASM plugin config: %w", err)
	}

	p := &Plugin{}

	if err := p.InitializeFromConfig(config); err != nil {
		return nil, fmt.Errorf("failed to initialize WASM plugin config: %w", err)
	}

	return p, nil
}

type record struct {
	K []byte            `json:"k,omitempty"`
	V []byte            `json:"v,omitempty"`
	H map[string]string `json:"h,omitempty"`
}

type Plugin struct {
	lock sync.Mutex

	r wazero.Runtime
	m wazapi.Module
	f wazapi.Function

	stdin  *bytes.Buffer
	stderr *bytes.Buffer
}

func (p *Plugin) InitializeFromConfig(config Config) error {
	ctx := context.Background()

	p.lock = sync.Mutex{}
	p.r = wazero.NewRuntime(ctx)

	// initialize WASI for stdin/out
	if _, err := wasi.NewBuilder(p.r).Instantiate(ctx); err != nil {
		return fmt.Errorf("failed to configurte WASI: %w", err)
	}

	moduleBytes, err := os.ReadFile(config.ModulePath)
	if err != nil {
		return fmt.Errorf("failed to read WASM module: %w", err)
	}

	p.stdin = bytes.NewBuffer(nil)
	p.stderr = bytes.NewBuffer(nil)

	mCfg := wazero.NewModuleConfig()
	mCfg = mCfg.WithStdin(p.stdin)
	mCfg = mCfg.WithStderr(p.stderr)

	if config.BindStdout {
		mCfg = mCfg.WithStdout(os.Stdout)
	}

	m, err := p.r.InstantiateWithConfig(ctx, moduleBytes, mCfg)
	if err != nil {
		return fmt.Errorf("failed to create WASM module: %w", err)
	}

	p.m = m
	p.f = p.m.ExportedFunction("produce")

	return nil
}

func (p *Plugin) Produce(ctx context.Context, k []byte, v []byte, h map[string]string) (*jrpc.ProduceResponse, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.stdin.Reset()
	p.stderr.Reset()

	data, err := json.Marshal(record{
		K: k,
		V: v,
		H: h,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to serialize WASM request: %w", err)
	}

	p.stdin.Write(data)
	ret, err := p.f.Call(ctx, uint64(len(data)))

	if err != nil {
		return nil, fmt.Errorf("failed to invoke WASM function: %w", err)
	}

	if len(ret) == 1 && ret[0] > 0 {
		err = p.extractError(ret[0])
		if err != nil {
			return nil, fmt.Errorf("failed to execute WASM function: %w", err)
		}
	}

	return &jrpc.ProduceResponse{
		Bytes:   0,
		Message: "",
	}, nil
}

func (p *Plugin) extractError(len uint64) error {
	if len == 0 {
		return nil
	}

	out := p.stderr.Bytes()
	if out == nil {
		return nil
	}

	return errors.New(string(out))
}
