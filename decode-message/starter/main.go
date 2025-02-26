package main

import (
	"context"
	"flag"
	"log"
	"strings"

	decodeMessage "github.com/josudoey/temporal-playground/decode-message"
	"go.temporal.io/sdk/client"
)

func main() {
	var input string
	flag.StringVar(&input, "i", "YXBwbGU=,YmFuYW5h", "strings of base64")
	flag.Parse()

	temporal, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporal.Close()

	ctx := context.Background()
	workflow, err := temporal.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        "decode_message_workflowID",
		TaskQueue: "decode-message",
	}, decodeMessage.Workflow, &decodeMessage.Input{
		Messages: strings.Split(input, ","),
	})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", workflow.GetID(), "RunID", workflow.GetRunID())

	var result decodeMessage.Result
	err = workflow.Get(ctx, &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
