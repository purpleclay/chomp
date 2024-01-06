/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
)

func TestWhile(t *testing.T) {
	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Happy New Year. Welcome 2024",
			rem:   " New Year. Welcome 2024",
			ext:   "Happy",
		},
		{
			name:  "Unicode",
			input: "あけましておめでとう。ようこそ 2024 年",
			rem:   "。ようこそ 2024 年",
			ext:   "あけましておめでとう",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rem, ext, err := chomp.While(chomp.IsLetter)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}

func TestWhileNot(t *testing.T) {
	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Happy New Year. Welcome 2024",
			rem:   "2024",
			ext:   "Happy New Year. Welcome ",
		},
		{
			name:  "Unicode",
			input: "あけましておめでとう。ようこそ 2024 年",
			rem:   "2024 年",
			ext:   "あけましておめでとう。ようこそ ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rem, ext, err := chomp.WhileNot(chomp.IsDigit)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}
