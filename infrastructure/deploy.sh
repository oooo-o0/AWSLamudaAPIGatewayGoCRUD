#!/bin/bash

set -e  # エラー発生時に停止

STACK_NAME="kaigo-insurance-system"
REGION="ap-northeast-1"  # 東京リージョンなど適宜変更
PROFILE="default"         # AWS CLIのプロファイル名。必要に応じて変更

# CloudFormationテンプレートディレクトリ
INFRA_DIR="infrastructure"

# デプロイ順序
# 1. IAMロールとポリシー（Lambdaなどの実行ロール）
# 2. DynamoDBテーブル
# 3. Lambda関数（IAMロール依存のため順序重要）
# 4. API Gateway
# 5. EventBridgeバッチスケジュール

echo "=== デプロイ開始 ==="

echo "1/5 IAMロールとポリシーのデプロイ"
aws cloudformation deploy \
  --template-file "${INFRA_DIR}/iam-roles.yaml" \
  --stack-name "${STACK_NAME}-iam" \
  --region "${REGION}" \
  --profile "${PROFILE}" \
  --capabilities CAPABILITY_NAMED_IAM

echo "2/5 DynamoDBテーブルのデプロイ"
aws cloudformation deploy \
  --template-file "${INFRA_DIR}/dynamodb.yaml" \
  --stack-name "${STACK_NAME}-dynamodb" \
  --region "${REGION}" \
  --profile "${PROFILE}" \
  --capabilities CAPABILITY_NAMED_IAM

echo "3/5 Lambda関数のデプロイ"
aws cloudformation deploy \
  --template-file "${INFRA_DIR}/lambda-functions.yaml" \
  --stack-name "${STACK_NAME}-lambda" \
  --region "${REGION}" \
  --profile "${PROFILE}" \
  --capabilities CAPABILITY_NAMED_IAM

echo "4/5 API Gatewayのデプロイ"
aws cloudformation deploy \
  --template-file "${INFRA_DIR}/api-gateway.yaml" \
  --stack-name "${STACK_NAME}-api" \
  --region "${REGION}" \
  --profile "${PROFILE}" \
  --capabilities CAPABILITY_NAMED_IAM

echo "5/5 EventBridgeバッチジョブスケジュールのデプロイ"
aws cloudformation deploy \
  --template-file "${INFRA_DIR}/batch-jobs.yaml" \
  --stack-name "${STACK_NAME}-batch-jobs" \
  --region "${REGION}" \
  --profile "${PROFILE}" \
  --capabilities CAPABILITY_NAMED_IAM

echo "=== デプロイ完了 ==="
