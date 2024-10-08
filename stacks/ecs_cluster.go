package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type EcsClusterStack struct {
	awscdk.Stack
	Cluster awsecs.Cluster
}

func NewEcsClusterStack(scope constructs.Construct, id string, vpc awsec2.Vpc) EcsClusterStack {
	stack := awscdk.NewStack(scope, &id, &awscdk.StackProps{
		Env: (*awscdk.Environment)(vpc.Env()),
	})

	cluster := awsecs.NewCluster(stack, jsii.String("ECS-Lab"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	return EcsClusterStack{
		Stack:   stack,
		Cluster: cluster,
	}
}
