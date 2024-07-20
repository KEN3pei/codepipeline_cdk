# Welcome to your CDK Go project!

This is a blank project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

### start

cdk init app --language go
go get

### codepipeline

1. codestarは手動で作成
2. cdk deploy '*' --profile <aws_account>

### 注意点

※ 現状トリガーが起動していない
