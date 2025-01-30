# Fiatconv 

Currency cli utility converter

Uses [https://exchangeratesapi.io](https://exchangeratesapi.io) as data source.

Based on the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) approaches.

Uses [contexts](https://golang.org/pkg/context/) where it possible.  

#### Arguments
- amount float
- src_symbol string
- dst_symbol string

#### Example
```shell
$ ./fiatconv 123.45 USD RUB
123.45 USD = 7894.12 RUB
```

#### Dependencies
Project holds all dependencies in `vendor` directory.

#### Testing

This project uses [mock](https://github.com/golang/mock) framework for testing.

To re-generate all mock files use:
```shell
$ go generate ./...
```

To run all tests:

```shell
$ go test -cover -v `go list ./...`
```