package main

import (
	"cdk/env"
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type CodePipelineProps struct {
	awscdk.StackProps
}

func NewCodePipelineStack(scope constructs.Construct, id string, props *CodePipelineProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// パイプラインv2の初期化
	pipeline := awscodepipeline.NewPipeline(stack, jsii.String(os.Getenv("PIPELINE_NAME")), &awscodepipeline.PipelineProps{
		PipelineName:     jsii.String(os.Getenv("PIPELINE_NAME")),
		CrossAccountKeys: jsii.Bool(false),
		PipelineType:     awscodepipeline.PipelineType_V2,
	})

	// Githubからソース取得
	sourceOutput := awscodepipeline.NewArtifact(jsii.String(os.Getenv("PROJECT") + "_Artifact"))
	sourceAction := awscodepipelineactions.NewCodeStarConnectionsSourceAction(&awscodepipelineactions.CodeStarConnectionsSourceActionProps{
		ActionName:         jsii.String(os.Getenv("PROJECT") + "_SourceAction"),
		Owner:              jsii.String(os.Getenv("SOURCE_OWNER")),
		Repo:               jsii.String(os.Getenv("SOURCE_REPOSITORY")),
		Output:             sourceOutput,
		Branch:             jsii.String(os.Getenv("SOURCE_BRANCH")),
		ConnectionArn:      jsii.String(os.Getenv("SOURCE_CONNECTION_ARN")),
		VariablesNamespace: env.GetNilOrStrEnv("SOURCE_NAME_SPACE"),
	})

	role := awsiam.NewRole(stack, jsii.String("CodeBuildTrustPolicy"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), nil),
	})
	// secret-managerアクセス用Role
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("secretsmanager:GetSecretValue"),
		},
		Resources: &[]*string{
			jsii.String(os.Getenv("SECRETS_MANAGER_ARN")),
		},
	}))
	// ECRアクセス・PUSH用Role
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("ecr:GetAuthorizationToken"),
			jsii.String("ecr:CompleteLayerUpload"),
			jsii.String("ecr:GetAuthorizationToken"),
			jsii.String("ecr:UploadLayerPart"),
			jsii.String("ecr:InitiateLayerUpload"),
			jsii.String("ecr:BatchCheckLayerAvailability"),
			jsii.String("ecr:PutImage"),
		},
		Resources: &[]*string{
			jsii.String(os.Getenv("ECR_REPOSITORY_ARN")),
		},
	}))

	// ビルドプロジェクト定義
	project := awscodebuild.NewProject(stack, jsii.String(os.Getenv("PROJECT")+"_BuildProject"), &awscodebuild.ProjectProps{
		BuildSpec: awscodebuild.BuildSpec_FromAsset(jsii.String("codebuild/buildspec.yml")),
		Role:      role,
	})
	// ビルドアクション定義
	buildAction := awscodepipelineactions.NewCodeBuildAction(&awscodepipelineactions.CodeBuildActionProps{
		ActionName: jsii.String(os.Getenv("PROJECT") + "_Action"),
		Project:    project,
		Input:      sourceOutput,
		EnvironmentVariables: &map[string]*awscodebuild.BuildEnvironmentVariable{
			"SECRETS_NAME": {
				Value: os.Getenv("SECRETS_MANAGER_ARN"),
			},
		},
	})

	// ソースステージ追加
	pipeline.AddStage(&awscodepipeline.StageOptions{
		StageName: jsii.String("Source"),
		Actions: &[]awscodepipeline.IAction{
			sourceAction,
		},
	})
	// ビルドステージ追加
	pipeline.AddStage(&awscodepipeline.StageOptions{
		StageName: jsii.String("Build"),
		Actions: &[]awscodepipeline.IAction{
			buildAction,
		},
	})

	// トリガーの設定
	pipeline.AddTrigger(&awscodepipeline.TriggerProps{
		ProviderType: awscodepipeline.ProviderType_CODE_STAR_SOURCE_CONNECTION,
		GitConfiguration: &awscodepipeline.GitConfiguration{
			SourceAction: sourceAction,
			PushFilter: &[]*awscodepipeline.GitPushFilter{
				{
					TagsIncludes: env.GetStringsEnv("GIT_PUSH_FILTER_INCLUDE_TAGS"),
				},
			},
		},
	})

	return stack
}

func main() {
	defer jsii.Close()

	loadEnv()

	app := awscdk.NewApp(nil)

	NewCodePipelineStack(app, "CodePipelineCDKStack", &CodePipelineProps{})

	app.Synth(nil)
}

func loadEnv() {
	err := godotenv.Load("env/" + "dev.env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}
