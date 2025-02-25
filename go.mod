module github.com/miyamo2/sqldav

go 1.21
toolchain go1.22.5

require (
	github.com/aws/aws-sdk-go-v2 v1.36.2
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.18.5
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.40.2
	github.com/google/go-cmp v0.6.0
	github.com/iancoleman/strcase v0.3.0
)

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.21 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)
