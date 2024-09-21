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

package random

import (
	"encoding/binary"
	"math/rand/v2"

	"github.com/google/uuid"
)

var JrSeed uint64
var ChaCha8 *rand.ChaCha8

var Random randomI

type randomI interface {
	IntN(n int) int
	Int64N(n int64) int64
	Float64() float64
	Float32() float32
	Uint64() uint64
	Shuffle(n int, swap func(i, j int))
}

// global random is a wrapper for the global random gen
type globalRandom struct {
}

func (r *globalRandom) IntN(n int) int {
	return rand.IntN(n) //nolint no need for a secure random generator
}
func (r *globalRandom) Int64N(n int64) int64 {
	return rand.Int64N(n) //nolint no need for a secure random generator
}

func (r *globalRandom) Float64() float64 {
	return rand.Float64() //nolint no need for a secure random generator
}
func (r *globalRandom) Float32() float32 {
	return rand.Float32() //nolint no need for a secure random generator
}
func (r *globalRandom) Shuffle(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap) //nolint no need for a secure random generator
}
func (r *globalRandom) Uint64() uint64 {
	return rand.Uint64() //nolint no need for a secure random generator
}
func SetRandom(seed int64) {

	if seed == -1 {
		// return a random/v2 object
		Random = &globalRandom{}
		return
	}

	JrSeed = uint64(seed) //nolint
	ChaCha8 = rand.NewChaCha8(CreateByteSeed(JrSeed))
	uuid.SetRand(ChaCha8)
	Random = rand.New(ChaCha8) //nolint no need for a secure random generator

}

func CreateByteSeed(seed uint64) [32]byte {
	b := make([]byte, 32)
	binary.LittleEndian.PutUint64(b, seed)
	binary.LittleEndian.PutUint64(b[8:], seed+1000)
	binary.LittleEndian.PutUint64(b[16:], seed+2000)
	binary.LittleEndian.PutUint64(b[24:], seed+3000)
	return [32]byte(b)
}
