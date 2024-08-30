//go:build tinygo.wasm

// tinygo build -o pkg/plugin/wasm/plugin_test_function.wasm -target=wasi pkg/plugin/wasm/plugin_test_function.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

const NoError = 0

type record struct {
	K []byte            `json:"k,omitempty"`
	V []byte            `json:"v,omitempty"`
	H map[string]string `json:"h,omitempty"`
}

//export produce
func _produce(size uint32) uint64 {
	b := make([]byte, size)

	_, err := io.ReadAtLeast(os.Stdin, b, int(size))
	if err != nil {
		return e(err)
	}

	in := record{}

	err = json.Unmarshal(b, &in)
	if err != nil {
		return e(err)
	}

	out := bytes.ToUpper(in.V)

	_, err = os.Stdout.Write(out)
	if err != nil {
		return e(err)
	}

	return NoError
}

func e(err error) uint64 {
	if err == nil {
		return NoError
	}

	_, _ = os.Stderr.WriteString(err.Error())
	return uint64(len(err.Error()))
}

func main() {}
