# Changelog

## 0.1.0 - 2024-07-13

### Initial Release🎉

implements the following `sql.Scanner`, `driver.Valuer`

- `sqldav.Set[string | int | float64 | []byte]`, the Defined Type of `[]string`, `[]int`, `[]float64`, `[][]byte`. Converted to `set` in DynamoDB.

- `sqldav.List`, the Defined Type of `[]interface{}`. Converted to `list` in DynamoDB.

- `sqldav.Map`, the Defined Type of `map[string]interface{}`. Converted to `map` in DynamoDB.

- `sqldav.TypedList[T]`, the Defined Type of `[]T`. Converted to `list` in DynamoDB.