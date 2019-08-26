# serviceexecutor
[![CircleCI](https://circleci.com/gh/cep21/serviceexecutor.svg)](https://circleci.com/gh/cep21/serviceexecutor)
[![GoDoc](https://godoc.org/github.com/cep21/serviceexecutor?status.svg)](https://godoc.org/github.com/cep21/serviceexecutor)
[![codecov](https://codecov.io/gh/cep21/serviceexecutor/branch/master/graph/badge.svg)](https://codecov.io/gh/cep21/serviceexecutor)

Serviceexecutor can manage long running service goroutines in go.

It assumes libraries follow best practices around not controlling concurrency for the user.  This abstraction
is intended to be used in one place: main.go.

# Usage

```go
    // TODO:
```

# Design Rational

TODO:

# Contributing

Contributions welcome!  Submit a pull request on github and make sure your code passes `make lint test`.  For
large changes, I strongly recommend [creating an issue](https://github.com/cep21/serviceexecutor/issues) on GitHub first to
confirm your change will be accepted before writing a lot of code.  GitHub issues are also recommended, at your discretion,
for smaller changes or questions.

# License

This library is licensed under the Apache 2.0 License.