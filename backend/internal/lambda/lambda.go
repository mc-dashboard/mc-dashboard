package lambda

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/rohanvsuri/minecraft-dashboard/internal/config"
)

type FunctionWrapper struct {
	LambdaClient *lambda.Client
}

func NewLambdaService(cfg *config.Config) *FunctionWrapper {
	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AwsRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AwsAccessKeyId, cfg.AwsSecretAccessKey, "")),
	)

	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
	}

	lambdaClient := lambda.NewFromConfig(awsConfig)
	return &FunctionWrapper{LambdaClient: lambdaClient}
}

func (f *FunctionWrapper) CallLambda(functionName string) *lambda.InvokeOutput {
	return f.Invoke(context.Background(), functionName, map[string]string{"key": "value"}, false)
}

func (wrapper FunctionWrapper) Invoke(ctx context.Context, functionName string, parameters any, getLog bool) *lambda.InvokeOutput {
	logType := types.LogTypeNone
	if getLog {
		logType = types.LogTypeTail
	}
	payload, err := json.Marshal(parameters)
	if err != nil {
		log.Panicf("Couldn't marshal parameters to JSON. Here's why %v\n", err)
	}
	invokeOutput, err := wrapper.LambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		LogType:      logType,
		Payload:      payload,
	})
	if err != nil {
		log.Panicf("Couldn't invoke function %v. Here's why: %v\n", functionName, err)
	}
	return invokeOutput
}
