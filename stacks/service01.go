package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapplicationautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Service01StackProps struct {
	awscdk.StackProps
	Cluster            awsecs.Cluster
	ProductEventsTopic awssns.Topic
}

func NewService01Stack(scope constructs.Construct, id string, props *Service01StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	rdsEndpoint := awscdk.Fn_ImportValue(jsii.String("rds-endpoint"))
	rdsUser := awscdk.Fn_ImportValue(jsii.String("rds-user"))
	rdsPass := awscdk.Fn_ImportValue(jsii.String("rds-pass"))
	dbTable := "products"

	envVariables := map[string]*string{
		"DB_DRIVER":       jsii.String("mysql"),
		"DB_HOST":         jsii.String(*rdsEndpoint),
		"DB_USER":         rdsUser,
		"DB_PASS":         rdsPass,
		"DB_NAME":         jsii.String("lab_aws"),
		"DB_TABLE":        &dbTable,
		"AWS_REGION":      jsii.String("sa-east-1"),
		"WEB_SERVER_PORT": jsii.String("8000"),
		// "AWS_SNS_TOPIC_PRODUCT_EVENTS_ARN": props.ProductEventsTopic
	}

	logGroup := awslogs.NewLogGroup(stack, jsii.String("Microservice-products-LG"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("Microservice-products"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	containerImage := awsecs.ContainerImage_FromRegistry(jsii.String("nayronferreiradev/microservice_products:v1.0.10"), nil)

	taskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("Task-Definition"), &awsecs.FargateTaskDefinitionProps{
		Cpu:            jsii.Number(256),
		MemoryLimitMiB: jsii.Number(512),
	})

	container := taskDefinition.AddContainer(jsii.String("aws_project01"), &awsecs.ContainerDefinitionOptions{
		Image: containerImage,
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			LogGroup:     logGroup,
			StreamPrefix: jsii.String("Microservice-products"),
		}),
		Environment: &envVariables,
	})

	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(8000),
	})

	service := awsecs.NewFargateService(stack, jsii.String("Microservice-products"), &awsecs.FargateServiceProps{
		Cluster:        props.Cluster,
		TaskDefinition: taskDefinition,
		DesiredCount:   jsii.Number(2),
		AssignPublicIp: jsii.Bool(false),
	})

	alb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("ALB-01"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		Vpc:            props.Cluster.Vpc(),
		InternetFacing: jsii.Bool(true),
	})

	listener := alb.AddListener(jsii.String("Listener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port: jsii.Number(80),
	})

	listener.AddTargets(jsii.String("ECS"), &awselasticloadbalancingv2.AddApplicationTargetsProps{
		Port:    jsii.Number(80),
		Targets: &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{service},
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Path:             jsii.String("/health"),
			Port:             jsii.String("8000"),
			HealthyHttpCodes: jsii.String("200"),
		},
	})

	scalableTaskCount := service.AutoScaleTaskCount(&awsapplicationautoscaling.EnableScalingProps{
		MinCapacity: jsii.Number(2),
		MaxCapacity: jsii.Number(3),
	})

	scalableTaskCount.ScaleOnCpuUtilization(jsii.String("Cpu-Scaling"), &awsecs.CpuUtilizationScalingProps{
		TargetUtilizationPercent: jsii.Number(50),
		ScaleInCooldown:          awscdk.Duration_Seconds(jsii.Number(60)),
		ScaleOutCooldown:         awscdk.Duration_Seconds(jsii.Number(60)),
	})

	scalableTaskCount.ScaleOnMemoryUtilization(jsii.String("Memory-Scaling"), &awsecs.MemoryUtilizationScalingProps{
		TargetUtilizationPercent: jsii.Number(50),
		ScaleInCooldown:          awscdk.Duration_Seconds(jsii.Number(60)),
		ScaleOutCooldown:         awscdk.Duration_Seconds(jsii.Number(60)),
	})

	// props.ProductEventsTopic.GrantPublish(service.TaskDefinition().TaskRole())

	return stack
}
