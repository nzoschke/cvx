// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package dynamodbiface provides an interface for the Amazon DynamoDB.
package dynamodbiface

import (
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBAPI is the interface type for dynamodb.DynamoDB.
type DynamoDBAPI interface {
	BatchGetItem(*dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error)

	BatchWriteItem(*dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)

	CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error)

	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)

	DeleteTable(*dynamodb.DeleteTableInput) (*dynamodb.DeleteTableOutput, error)

	DescribeTable(*dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error)

	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)

	ListTables(*dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error)

	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)

	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)

	Scan(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error)

	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)

	UpdateTable(*dynamodb.UpdateTableInput) (*dynamodb.UpdateTableOutput, error)
}