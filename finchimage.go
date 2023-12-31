package main

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FinchContainerStackProps struct {
	awscdk.StackProps
}

func FinchContainerStack(scope constructs.Construct, id string, props *FinchContainerStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create new container image
	dir, _ := os.Getwd()
	dirImage := filepath.Join(dir, "imageassets")
	buildArgs := false

	asset := awsecrassets.NewDockerImageAsset(stack, jsii.String("MyBuildImage"), &awsecrassets.DockerImageAssetProps{
		Directory: &dirImage,
		Invalidation: &awsecrassets.DockerImageAssetInvalidationOptions{
			BuildArgs: &buildArgs,
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("Image built"), &awscdk.CfnOutputProps{Value: asset.ImageUri()})

	// Create VPC and Cluster
	vpc := awsec2.NewVpc(stack, jsii.String("ALBFargoVpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})

	cluster := awsecs.NewCluster(stack, jsii.String("ALBFargoECSCluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	res := awsecspatterns.NewApplicationLoadBalancedFargateService(stack, jsii.String("ALBFargoService"), &awsecspatterns.ApplicationLoadBalancedFargateServiceProps{
		Cluster: cluster,
		TaskImageOptions: &awsecspatterns.ApplicationLoadBalancedTaskImageOptions{
			Image:         awsecs.AssetImage_FromDockerImageAsset(asset),
			ContainerPort: jsii.Number(80),
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("LoadBalancerDNS"), &awscdk.CfnOutputProps{Value: res.LoadBalancer().LoadBalancerDnsName()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	FinchContainerStack(app, "FinchContainerStack", &FinchContainerStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
