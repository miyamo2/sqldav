module github.com/miyamo2/sqldav

go 1.21
toolchain go1.24.1

require (
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.15
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.42.0
	github.com/google/go-cmp v0.6.0
	github.com/iancoleman/strcase v0.3.0
)

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.5 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)
