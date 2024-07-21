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

### 事前準備

- codestar-connectionn（code-conneciton）でGitHubとの接続を作成しておく
- secret-manager作成
- CDK Bootstrap実行しておく（[bootstrap](https://docs.aws.amazon.com/ja_jp/cdk/v2/guide/bootstrapping.html)）

### 実行コマンド

$ cdk deploy '*' --profile <aws_account>

### トラブルシュート

※ トリガーが起動しない
- GitHub側のAWS Connector for GitHubでトリガー対象リポジトリが許可されていない可能性がある。

