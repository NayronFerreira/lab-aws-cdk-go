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

	// envVariables := map[string]*string{
	// 	"SPRING_DATASOURCE_URL":            jsii.String("jdbc:mysql://" + awscdk.Fn_ImportValue(jsii.String("rds-endpoint")) + ":3306/aws_project01?createDatabaseIfNotExist=true"),
	// 	"SPRING_DATASOURCE_USERNAME":       awscdk.Fn_ImportValue(jsii.String("rds-user")),
	// 	"SPRING_DATASOURCE_PASSWORD":       awscdk.Fn_ImportValue(jsii.String("rds-pass")),
	// 	"AWS_REGION":                       jsii.String("us-east-1"),
	// 	"AWS_SNS_TOPIC_PRODUCT_EVENTS_ARN": props.ProductEventsTopic.TopicArn(),
	// }

	logGroup := awslogs.NewLogGroup(stack, jsii.String("Service-01-LogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("Service-01"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	containerImage := awsecs.ContainerImage_FromRegistry(jsii.String("nayronferreiradev/microservice_products:v1.0.1"), nil)

	taskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("Task-Definition"), &awsecs.FargateTaskDefinitionProps{
		Cpu:            jsii.Number(256),
		MemoryLimitMiB: jsii.Number(512),
	})

	container := taskDefinition.AddContainer(jsii.String("aws_project01"), &awsecs.ContainerDefinitionOptions{
		Image: containerImage,
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			LogGroup:     logGroup,
			StreamPrefix: jsii.String("Service-01"),
		}),
		// Environment: props.Cluster.Vpc(),
	})

	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(8000),
	})

	service := awsecs.NewFargateService(stack, jsii.String("Service-01"), &awsecs.FargateServiceProps{
		Cluster:        props.Cluster,
		TaskDefinition: taskDefinition,
		DesiredCount:   jsii.Number(2),
		AssignPublicIp: jsii.Bool(true),
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
