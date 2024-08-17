package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RdsStackProps struct {
	awscdk.StackProps
	Vpc              awsec2.Vpc
	DatabaseUser     *string
	DatabasePassword *string
}

func NewRdsStack(scope constructs.Construct, id string, props *RdsStackProps) awscdk.Stack {

	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Parâmetros do CloudFormation para o usuário e senha do banco de dados
	databaseUser := awscdk.NewCfnParameter(stack, jsii.String("databaseUser"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("RDS Database user"),
		Default:     props.DatabaseUser,
	})

	databasePass := awscdk.NewCfnParameter(stack, jsii.String("databasePass"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("RDS Database password"),
		Default:     props.DatabasePassword,
	})

	// Grupo de segurança
	securityGroup := awsec2.SecurityGroup_FromSecurityGroupId(stack,
		jsii.String("SecurityGroup"),
		props.Vpc.VpcDefaultSecurityGroup(), nil)

	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow MySQL access"), jsii.Bool(false))

	// Definindo as propriedades da instância RDS
	instanceProps := &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Mysql(&awsrds.MySqlInstanceEngineProps{
			Version: awsrds.MysqlEngineVersion_VER_8_0(),
		}),
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE3, awsec2.InstanceSize_MICRO),
		Vpc:          props.Vpc,
		Credentials: awsrds.Credentials_FromUsername(databaseUser.ValueAsString(), &awsrds.CredentialsFromUsernameOptions{
			Password: awscdk.SecretValue_UnsafePlainText(databasePass.ValueAsString()),
		}),
		MultiAz:          jsii.Bool(false),
		AllocatedStorage: jsii.Number(10),
		SecurityGroups:   &[]awsec2.ISecurityGroup{securityGroup},
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS, // Sub-rede privada com acesso à internet via NatGateway
		},
	}

	// Criando a instância RDS
	databaseInstance := awsrds.NewDatabaseInstance(stack, jsii.String("RDS-MySQL"), instanceProps)

	// Saídas do CloudFormation
	awscdk.NewCfnOutput(stack, jsii.String("rds-endpoint"), &awscdk.CfnOutputProps{
		ExportName:  jsii.String("rds-endpoint"),
		Description: jsii.String("RDS Endpoint"),
		Value:       databaseInstance.DbInstanceEndpointAddress(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("rds-password"), &awscdk.CfnOutputProps{
		ExportName:  jsii.String("rds-pass"),
		Description: jsii.String("RDS Password"),
		Value:       databasePass.ValueAsString(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("rds-user"), &awscdk.CfnOutputProps{
		ExportName:  jsii.String("rds-user"),
		Description: jsii.String("RDS Username"),
		Value:       databaseUser.ValueAsString(),
	})

	return stack
}
