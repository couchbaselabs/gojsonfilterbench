# Go JSON Filter Benchmarks
This repository contains benchmark comparisons between various possible
implementations of performing filtering upon a stream of JSON documents.

# Results
Here is an example run on my 2017 Macbook Pro 2.9Ghz:
```
> go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/couchbaselabs/gojsonfilterbench
BenchmarkJsonPathWithGoValuate-8   	     500	   2686110 ns/op	 359.87 MB/s
BenchmarkJsonSM-8                  	    1000	   1606162 ns/op	 619.65 MB/s
BenchmarkJsonSMSlowMatcher-8       	     100	  17878969 ns/op	  53.56 MB/s
PASS
ok  	github.com/couchbaselabs/gojsonfilterbench	5.759s
```
