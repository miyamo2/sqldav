# Changelog

## 0.2.1 - 2024-10-06

### ⬆️ Upgrading dependencies

-  `github.com/aws/aws-sdk-go-v2` from 1.30.4 to 1.31.0
-  `github.com/aws/aws-sdk-go-v2/service/dynamodb` from v1.34.3 to v1.36.0
-  `github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue` from v1.14.9 to v1.15.10

## 0.2.0 - 2024-07-13

### ✨ New Features

- `Set`, `List`, `Map` and `TypedList` now implements `gorm.io/gorm.schema.GormDataTypeInterface`.

## 0.1.1 - 2024-07-13

### Bug Fix🐛

#### `driver.Valuer` implementations

Receiver is now a physical value instead of a pointer.
This ensures that type switches and type assertions work properly.

## 0.1.0 - 2024-07-13

### Initial Release🎉

implements the following `sql.Scanner`, `driver.Valuer`

- `sqldav.Set[string | int | float64 | []byte]`, the Defined Type of `[]string`, `[]int`, `[]float64`, `[][]byte`. Converted to `set` in DynamoDB.

- `sqldav.List`, the Defined Type of `[]interface{}`. Converted to `list` in DynamoDB.

- `sqldav.Map`, the Defined Type of `map[string]interface{}`. Converted to `map` in DynamoDB.

- `sqldav.TypedList[T]`, the Defined Type of `[]T`. Converted to `list` in DynamoDB.
