### Steps to run this sample:
1) Run a [Temporal service](https://github.com/temporalio/samples-go/tree/main/#how-to-use).
2) Run the following command to start the worker
```
go run decode-message/worker/main.go
```
3) Run the following command to start the example
```
go run decode-message/starter/main.go -i YmFuYW5h,YmFuYW5hCg==
```
4) Run the following command to query the example
go run decode-message/start/main.go -r <RunID>
