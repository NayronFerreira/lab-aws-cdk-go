package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VpcStack struct {
	awscdk.Stack
	Vpc awsec2.Vpc
}

func NewVpcStack(scope constructs.Construct, id string, env *awscdk.Environment) VpcStack {
	stack := awscdk.NewStack(scope, &id, &awscdk.StackProps{
		Env: env,
	})

	vpc := awsec2.NewVpc(stack, jsii.String("Vpc"), &awsec2.VpcProps{
		MaxAzs:      jsii.Number(2),
		NatGateways: jsii.Number(0),
	})

	return VpcStack{
		Stack: stack,
		Vpc:   vpc,
	}
}
