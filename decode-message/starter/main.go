package main

import (
	"context"
	"flag"
	"log"
	"reflect"
	"strings"
	"time"

	decodeMessage "github.com/josudoey/temporal-playground/decode-message"
	"go.temporal.io/sdk/client"
)

func main() {
	var workflowID, input string
	flag.StringVar(&workflowID, "w", "decode_message_workflowID", "WorkflowID.")
	flag.StringVar(&input, "i", "YXBwbGU=,YmFuYW5h", "strings of base64")
	flag.Parse()

	temporal, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporal.Close()

	ctx := context.Background()
	workflow, err := temporal.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: decodeMessage.TaskQueueName,
	}, decodeMessage.Workflow, &decodeMessage.Input{
		Messages: strings.Split(input, ","),
	})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", workflow.GetID(), "RunID", workflow.GetRunID())

	for {
		var progress decodeMessage.Progress
		resp, err := temporal.QueryWorkflow(ctx, workflow.GetID(), workflow.GetRunID(), reflect.TypeOf(progress).Name())
		if err != nil {
			log.Fatalln("Unable to query workflow", err)
		}
		if err := resp.Get(&progress); err != nil {
			log.Fatalln("Unable to decode query result", err)
		}
		if progress.CurrentCount == progress.Total {
			break
		}
		log.Printf("progress: %+v", progress)
		time.Sleep(time.Second)
	}

	var result decodeMessage.Result
	err = workflow.Get(ctx, &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
