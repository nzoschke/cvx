// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package dynamodbiface_test

import (
	"testing"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	assert.Implements(t, (*dynamodbiface.DynamoDBAPI)(nil), dynamodb.New(nil))
}