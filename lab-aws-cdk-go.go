package main

import (
	"os"

	"github.com/NayronFerreira/lab-aws-cdk-go/stacks"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

func main() {
	defer jsii.Close()

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	app := awscdk.NewApp(nil)

	vpcStack := stacks.NewVpcStack(app, "VPC", env())

	ecsClusterStack := stacks.NewEcsClusterStack(app, "Cluster", vpcStack.Vpc)
	ecsClusterStack.AddDependency(vpcStack.Stack, nil)

	rdsStack := stacks.NewRdsStack(app, "RDS", &stacks.RdsStackProps{
		StackProps: awscdk.StackProps{
			Env: (*awscdk.Environment)(vpcStack.Vpc.Env()),
		},
		Vpc:              vpcStack.Vpc,
		DatabaseUser:     jsii.String(os.Getenv("RDS_USER")),
		DatabasePassword: jsii.String(os.Getenv("RDS_PASS")),
	})
	rdsStack.AddDependency(vpcStack.Stack, nil)

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
	service01Stack.AddDependency(rdsStack, nil)

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
