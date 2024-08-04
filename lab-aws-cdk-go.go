package main

import (
	"github.com/NayronFerreira/lab-aws-cdk-go/stacks"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	vpcStack := stacks.NewVpcStack(app, "VPC", env())

	ecsClusterStack := stacks.NewEcsClusterStack(app, "ECS-Cluster", vpcStack.Vpc)
	ecsClusterStack.AddDependency(vpcStack.Stack, nil)

	// Cria a stack do Service01, passando o cluster ECS criado
	service01Props := &stacks.Service01StackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		Cluster: ecsClusterStack.Cluster,
		// ProductEventsTopic: productEventsTopic,
	}

	service01Stack := stacks.NewService01Stack(app, "Service-01", service01Props)
	service01Stack.AddDependency(ecsClusterStack.Stack, nil)

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
