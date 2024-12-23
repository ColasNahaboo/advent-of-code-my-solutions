// Boyer-Moore fast search of a substring
// From https://github.com/sarpdag/boyermoore (MIT license at end)

// Note: in case of multiple search of the same substring,
// re-use the result of BMStringCalculateSlideTable in call to
// BMStringIndexWithTable

// Index returns the first index substr found in the s.
// function should return same result as `strings.Index` function
func BMStringIndex(s string, substr string) int {
	d := BMStringCalculateSlideTable(substr)
	return BMStringIndexWithTable(&d, s, substr)
}

// IndexWithTable returns the first index substr found in the s.
// It needs the slide information of substr
func BMStringIndexWithTable(d *[256]int, s string, substr string) int {
	lsub := len(substr)
	ls := len(s)
	// fmt.Println(ls, lsub)
	switch {
	case lsub == 0:
		return 0
	case lsub > ls:
		return -1
	case lsub == ls:
		if s == substr {
			return 0
		}
		return -1
	}

	i := 0
	for i+lsub-1 < ls {
		j := lsub - 1
		for ; j >= 0 && s[i+j] == substr[j]; j-- {
		}
		if j < 0 {
			return i
		}

		slid := j - d[s[i+j]]
		if slid < 1 {
			slid = 1
		}
		i += slid
	}
	return -1
}

// CalculateSlideTable builds sliding amount per each unique byte in the substring
func BMStringCalculateSlideTable(substr string) [256]int {
	var d [256]int
	for i := 0; i < 256; i++ {
		d[i] = -1
	}
	for i := 0; i < len(substr); i++ {
		d[substr[i]] = i
	}
	return d
}


// MIT License
// 
// Copyright (c) 2021 Sarp Dağ Demirel
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
