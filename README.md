[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/sanjibdevnathlabs/gosafejson)
[![codecov](https://codecov.io/gh/sanjibdevnathlabs/gosafejson/branch/master/graph/badge.svg)](https://codecov.io/gh/sanjibdevnathlabs/gosafejson)
[![rcard](https://goreportcard.com/badge/github.com/sanjibdevnathlabs/gosafejson)](https://goreportcard.com/report/github.com/sanjibdevnathlabs/gosafejson)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/sanjibdevnathlabs/gosafejson/blob/master/LICENSE)

**Note:** This repository (`gosafejson`) is a fork of the original [json-iterator/go](https://github.com/json-iterator/go).

A high-performance 100% compatible drop-in replacement of "encoding/json"

# Benchmark

![benchmark](http://jsoniter.com/benchmarks/go-benchmark.png)

Source code: https://github.com/json-iterator/go-benchmark/blob/master/src/github.com/json-iterator/go-benchmark/benchmark_medium_payload_test.go

Raw Result (easyjson requires static code generation)

|                 | ns/op       | allocation bytes | allocation times |
| --------------- | ----------- | ---------------- | ---------------- |
| std decode      | 35510 ns/op | 1960 B/op        | 99 allocs/op     |
| easyjson decode | 8499 ns/op  | 160 B/op         | 4 allocs/op      |
| jsoniter decode | 5623 ns/op  | 160 B/op         | 3 allocs/op      |
| std encode      | 2213 ns/op  | 712 B/op         | 5 allocs/op      |
| easyjson encode | 883 ns/op   | 576 B/op         | 3 allocs/op      |
| jsoniter encode | 837 ns/op   | 384 B/op         | 4 allocs/op      |

Always benchmark with your own workload.
The result depends heavily on the data input.

# Usage

100% compatibility with standard lib

Replace

```go
import "encoding/json"
json.Marshal(&data)
```

with

```go
import "github.com/sanjibdevnathlabs/gosafejson"

var json = gosafejson.ConfigCompatibleWithStandardLibrary
json.Marshal(&data)
```

Replace

```go
import "encoding/json"
json.Unmarshal(input, &data)
```

with

```go
import "github.com/sanjibdevnathlabs/gosafejson"

var json = gosafejson.ConfigCompatibleWithStandardLibrary
json.Unmarshal(input, &data)
```

# Safe Unmarshalling

`gosafejson` provides a safe unmarshalling mode that continues processing JSON even when type mismatches occur. This is useful when you want to extract as much valid data as possible from a JSON document, rather than failing completely on the first error.

With standard unmarshalling (both in Go's standard library and in default mode), a type mismatch will cause the entire unmarshalling process to fail:

```go
type Details struct {
    ID   *string                `json:"id"`
    Data map[string]interface{} `json:"data"`
    Name *string                `json:"name"`
}

// Type mismatch: "data" is an array, not a map
jsonStr := `{"id":"12345", "data": [{"a":"b"}, {"c":"d"}], "name": "Random"}`

var details Details
err := gosafejson.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(jsonStr), &details)
// err will contain an error, and details will be incomplete
```

With safe unmarshalling, the process continues after errors, and you get a composite error with all the issues encountered:

```go
var details Details
err := gosafejson.ConfigSafe.Unmarshal([]byte(jsonStr), &details)
// err will be of type *CompositeError with details on all type mismatches
// Fields before the error will be properly unmarshalled (like ID in this example)
```

This is particularly useful when:
- Working with inconsistent or evolving APIs
- Processing data from third-party sources that may have inconsistencies
- You need to extract partial data even when the full document has errors

[More documentation](http://jsoniter.com/migrate-from-go-std.html)

# How to get

```
go get github.com/sanjibdevnathlabs/gosafejson
```
