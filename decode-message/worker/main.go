package main

import (
	"log"

	decodeMessage "github.com/josudoey/temporal-playground/decode-message"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	temporal, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporal.Close()

	w := worker.New(temporal, decodeMessage.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(decodeMessage.Workflow)
	w.RegisterActivity(decodeMessage.DecodeMessage)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
