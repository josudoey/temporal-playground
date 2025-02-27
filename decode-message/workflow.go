package decodemessage

import (
	"context"
	"encoding/base64"
	"reflect"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueueName = "decode-message-tq"
)

type Input struct {
	Messages []string
}

type Progress struct {
	Total        int
	CurrentCount int
	SuccessCount int
	FailedCount  int
}

type Result struct {
	Total        int
	SuccessCount int
	FailedCount  int
}

func Workflow(ctx workflow.Context, input *Input) (*Result, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	})

	logger := workflow.GetLogger(ctx)
	progress := Progress{Total: len(input.Messages)}
	err := workflow.SetQueryHandler(ctx, reflect.TypeOf(progress).Name(), func(input []byte) (*Progress, error) {
		return &progress, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
		return nil, err
	}

	for _, message := range input.Messages {
		var answer string

		progress.CurrentCount += 1
		if err := workflow.ExecuteActivity(ctx, DecodeMessage, message).Get(ctx, &answer); err != nil {
			logger.Error("decode message failed", "err", err)
			progress.FailedCount += 1
			continue
		}
		progress.SuccessCount += 1
	}

	return &Result{
		Total:        progress.Total,
		SuccessCount: progress.SuccessCount,
		FailedCount:  progress.FailedCount,
	}, nil
}

func DecodeMessage(ctx context.Context, message string) (string, error) {
	logger := activity.GetLogger(ctx)
	time.Sleep(1 * time.Second)
	result, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}
	logger.Info("encode message", "result", result)
	return string(result), nil
}
