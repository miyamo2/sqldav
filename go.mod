module github.com/miyamo2/sqldav

go 1.21
toolchain go1.22.5

require (
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.18.6
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.41.0
	github.com/google/go-cmp v0.6.0
	github.com/iancoleman/strcase v0.3.0
)

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.25.0 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)
