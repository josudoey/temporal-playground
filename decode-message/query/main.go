package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/history/v1"
	"go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func main() {
	var workflowID, runID string
	flag.StringVar(&workflowID, "w", "decode_message_workflowID", "WorkflowID.")
	flag.StringVar(&runID, "r", "", "RunID")
	flag.Parse()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "WorkflowId=" + strconv.Quote(workflowID)
	if runID != "" {
		query += " and RunId=" + strconv.Quote(runID)
	}

	// query ref https://docs.temporal.io/search-attribute#default-search-attribute
	res, err := c.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: 30,
		Query:    query,
	})
	if err != nil {
		log.Fatalln("list workflow error", err)
	}

	for _, workflowExecutionInfo := range res.GetExecutions() {
		printWorkflowExecutionInfo(workflowExecutionInfo)
		execution := workflowExecutionInfo.GetExecution()
		iter := c.GetWorkflowHistory(ctx, execution.GetWorkflowId(), execution.GetRunId(), false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
		for iter.HasNext() {
			e, err := iter.Next()
			if err != nil {
				log.Fatalln("get history error", err)
			}

			printEvent(e)
		}
	}
}

func printWorkflowExecutionInfo(info *workflow.WorkflowExecutionInfo) {
	fmt.Printf("%s ", info.GetExecutionTime().AsTime().Format(time.RFC3339Nano))
	fmt.Printf("workflow duration: %s ", info.GetExecutionDuration().AsDuration())
	fmt.Printf("WorkflowId: %s RunId: %s \n", info.GetExecution().GetWorkflowId(), info.GetExecution().GetRunId())
}

func printEventTime(e *history.HistoryEvent) {
	fmt.Printf("%s ", e.GetEventTime().AsTime().Format(time.RFC3339Nano))
}

func printEvent(e *history.HistoryEvent) {
	switch e.GetEventType() {
	case enums.EVENT_TYPE_WORKFLOW_EXECUTION_STARTED:
		printEventTime(e)
		fmt.Printf("workflow execution started input: %s\n", e.GetWorkflowExecutionStartedEventAttributes().GetInput().GetPayloads()[0].GetData())
	case enums.EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED:
		printEventTime(e)
		fmt.Printf("workflow execution completed result: %s\n", e.GetWorkflowExecutionCompletedEventAttributes().GetResult().GetPayloads()[0].GetData())
	case enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
		printEventTime(e)
		fmt.Printf("activity scheduled scheduled: %s\n", e.GetActivityTaskScheduledEventAttributes().GetInput().GetPayloads()[0].GetData())
	case enums.EVENT_TYPE_ACTIVITY_TASK_STARTED:
		printEventTime(e)
		fmt.Printf("activity started attempt: %d\n", e.GetActivityTaskStartedEventAttributes().GetAttempt())
	case enums.EVENT_TYPE_ACTIVITY_TASK_COMPLETED:
		printEventTime(e)
		fmt.Printf("activity completed result: %s\n", e.GetActivityTaskCompletedEventAttributes().GetResult().GetPayloads()[0].GetData())
	case enums.EVENT_TYPE_ACTIVITY_TASK_FAILED:
		printEventTime(e)
		fmt.Printf("activity failed message: %s\n", e.GetActivityTaskFailedEventAttributes().GetFailure().GetMessage())
	}
}
