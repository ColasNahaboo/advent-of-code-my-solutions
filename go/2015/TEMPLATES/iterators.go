// Fastest iterator as of Go 1.17
// From https://ewencp.org/blog/golang-iterators/index.html

// goos: linux
// goarch: amd64
// pkg: github.com/ewencp/golang-iterators-benchmark
// cpu: 11th Gen Intel(R) Core(TM) i5-11500 @ 2.70GHz
// BenchmarkIntsCallbackIterator-12         787           1485837 ns/op
// BenchmarkDataCallbackIterator-12         726           1631669 ns/op
// BenchmarkIntsChannelIterator-12            9         124790870 ns/op
// BenchmarkDataChannelIterator-12            9         125658045 ns/op
// BenchmarkIntsBufferedChannelIterator-12   33          36467180 ns/op
// BenchmarkDataBufferedChannelIterator-12   34          35227674 ns/op
// BenchmarkIntsClosureIterator-12         2653            450429 ns/op
// BenchmarkDataClosureIterator-12         1132           1075460 ns/op
// BenchmarkIntStatefulIterator-12         3553            339929 ns/op
// BenchmarkDataStatefulIterator-12        1104           1211424 ns/op
// BenchmarkIntStatefulIteratorInterface-12 596           1910526 ns/op
// BenchmarkDataStatefulIteratorInterface-12768           1552404 ns/op


// The iterator: builds and return a iterator function
// replace the "int_data" lines with the implentation for your case
func IntClosureIterator() (func() (int, bool), bool) {
    var idx int = 0
    var data_len = len(int_data)
    return func() (int, bool) {
        prev_idx := idx
        idx++
        return int_data[prev_idx], (idx < data_len)
    }, (idx < data_len)
}

// calling it:

var sum, val int = 0, 0
for it, has_next := IntClosureIterator(); has_next; val, has_next = it() {
    sum += val
}

// same for structs:
type Data struct {
	foo int
	bar *Data
}
var struct_data []*Data = make([]*Data, NumItems)

func DataClosureIterator() (func() (int, bool), bool) {
	var idx int = 0
	var data_len = len(struct_data)
	return func() (int, bool) {
		prev_idx := idx
		idx++
		return struct_data[prev_idx].foo, (idx < data_len)
	}, (idx < data_len)
}
