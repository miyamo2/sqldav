# sqldav

[![Go Reference](https://pkg.go.dev/badge/github.com/miyamo2/sqldav.svg)](https://pkg.go.dev/github.com/miyamo2/sqldav)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/miyamo2/sqldav?logo=go)](https://img.shields.io/github/go-mod/go-version/miyamo2/sqldav?logo=go)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/miyamo2/sqldav)](https://img.shields.io/github/v/release/miyamo2/sqldav)
[![codecov](https://codecov.io/gh/miyamo2/sqldav/graph/badge.svg?token=J0RZ235JE4)](https://codecov.io/gh/miyamo2/sqldav)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/sqldav)](https://goreportcard.com/report/github.com/miyamo2/sqldav)
[![GitHub License](https://img.shields.io/github/license/miyamo2/sqldav?&color=blue)](https://img.shields.io/github/license/miyamo2/sqldav?&color=blue)

"sql.Scanner"/"driver.Valuer" for DynamoDB PartiQL (And its Tooltip)

## Types

sqldav implements the following `sql.Scanner`, `driver.Valuer`

- `sqldav.Set[string | int | float64 | []byte]`, the Defined Type of `[]string`, `[]int`, `[]float64`, `[][]byte`. Converted to `set` in DynamoDB.

- `sqldav.List`, the Defined Type of `[]interface{}`. Converted to `list` in DynamoDB.

- `sqldav.Map`, the Defined Type of `map[string]interface{}`. Converted to `map` in DynamoDB.

- `sqldav.TypedList[T]`, the Defined Type of `[]T`. Converted to `list` in DynamoDB.

## Contributing

Feel free to open a PR or an Issue.

However, you must promise to follow our [Code of Conduct](https://github.com/miyamo2/sqldav/blob/main/CODE_OF_CONDUCT.md).

## License

**sqldav** released under the [MIT License](https://github.com/miyamo2/sqldav/blob/main/LICENSE)