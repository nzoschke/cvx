// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package cloudformation

import (
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/internal/protocol/query"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/internal/signer/v4"
)

// AWS CloudFormation enables you to create and manage AWS infrastructure deployments
// predictably and repeatedly. AWS CloudFormation helps you leverage AWS products
// such as Amazon EC2, EBS, Amazon SNS, ELB, and Auto Scaling to build highly-reliable,
// highly scalable, cost effective applications without worrying about creating
// and configuring the underlying AWS infrastructure.
//
// With AWS CloudFormation, you declare all of your resources and dependencies
// in a template file. The template defines a collection of resources as a single
// unit called a stack. AWS CloudFormation creates and deletes all member resources
// of the stack together and manages all dependencies between the resources
// for you.
//
// For more information about this product, go to the CloudFormation Product
// Page (http://aws.amazon.com/cloudformation/).
//
// Amazon CloudFormation makes use of other AWS products. If you need additional
// technical information about a specific AWS product, you can find the product's
// technical documentation at http://aws.amazon.com/documentation/ (http://aws.amazon.com/documentation/).
type CloudFormation struct {
	*aws.Service
}

// Used for custom service initialization logic
var initService func(*aws.Service)

// Used for custom request initialization logic
var initRequest func(*aws.Request)

// New returns a new CloudFormation client.
func New(config *aws.Config) *CloudFormation {
	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "cloudformation",
		APIVersion:  "2010-05-15",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(query.Build)
	service.Handlers.Unmarshal.PushBack(query.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(query.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(query.UnmarshalError)

	// Run custom service initialization if present
	if initService != nil {
		initService(service)
	}

	return &CloudFormation{service}
}

// newRequest creates a new request for a CloudFormation operation and runs any
// custom request initialization.
func (c *CloudFormation) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
