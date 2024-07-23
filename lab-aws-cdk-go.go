package main

import (
	"github.com/NayronFerreira/lab-aws-cdk-go/stacks"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	// Cria a stack da VPC
	vpcStack := stacks.NewVpcStack(app, "VPC", env())

	// Cria a stack do Cluster ECS, passando a VPC criada
	ecsClusterStack := stacks.NewEcsClusterStack(app, "ECS-Cluster", vpcStack.Vpc, env())
	ecsClusterStack.AddDependency(vpcStack.Stack, nil)

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
