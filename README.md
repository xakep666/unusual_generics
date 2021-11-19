Unusual Generics
================

[![Run Tests](https://github.com/xakep666/unusual_generics/actions/workflows/testing.yml/badge.svg)](https://github.com/xakep666/unusual_generics/actions/workflows/testing.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/xakep666/unusual_generics.svg)](https://pkg.go.dev/github.com/xakep666/unusual_generics)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[Type parameters](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) or Generics
in Go designed to reduce boilerplate for container data types like lists, graphs, etc. and functions like map, filter, reduce...

But it's possible to use them in other (I've named them 'unusual') cases. I'm collecting such cases in this repository.
What's inside:
* [Type to emulate JS 'undefined' for JSON](json_undefined.go)
* [Type to deal with non-standard time formats in JSON/XML/etc.](time_format.go)
* [Function to get pointer from literal in one line](ptr.go)
* [Generic version of x/sync/singleflight.Group](singleflight.go)
* [Reflect(lite)-free versions of errors.Is and errors.As with upgrade](errors.go)

Feel free to open issue or pull request to add new one.

Note that until 1.18 release documentation on 'pkg.go.dev' will not be rendered.
