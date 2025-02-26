package decodemessage

import (
	"context"
	"encoding/base64"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Input struct {
	Messages []string
}

type Result struct {
	Total        int
	SuccessCount int
	FailedCount  int
}

func Workflow(ctx workflow.Context, input *Input) (*Result, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 200 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	})

	logger := workflow.GetLogger(ctx)
	result := &Result{
		Total: len(input.Messages),
	}
	for _, message := range input.Messages {
		var answer string
		if err := workflow.ExecuteActivity(ctx, DecodeMessage, message).Get(ctx, &answer); err != nil {
			logger.Error("decode message failed", "err", err)
			result.FailedCount += 1
			continue
		}
		result.SuccessCount += 1
	}

	return result, nil
}

func DecodeMessage(ctx context.Context, message string) (string, error) {
	logger := activity.GetLogger(ctx)
	result, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}
	logger.Info("encode message", "result", result)
	return string(result), nil
}
