# Loki Client

The Loki Client is a Go library for sending logs to [Loki](https://github.com/grafana/loki). It provides an efficient way to batch, compress, and push logs over HTTP using snappy-compressed protobufs.

## Features

- Batch logs for efficient transmission.
- Snappy compression for reduced payload size.
- Configurable backoff and retry mechanism for handling transient errors.
- Support for external labels to enrich log streams.
- Implements the Loki Push API.

## Installation

To use the Loki Client in your Go project, add it as a dependency:

```sh
go get github.com/livepeer/loki-client
```

## Usage

### Creating a Client

You can create a new Loki client using the `NewWithDefaults` or `New` functions.

```go
package main

import (
	"log"
	"time"

	"github.com/livepeer/loki-client/client"
	"github.com/livepeer/loki-client/model"
)

func main() {
	logger := func(v ...interface{}) {
		log.Println(v...)
	}

	externalLabels := model.LabelSet{
		"app": "my-app",
	}

	client, err := client.NewWithDefaults("http://localhost:3100/loki/api/v1/push", externalLabels, logger)
	if err != nil {
		log.Fatalf("Failed to create Loki client: %v", err)
	}
	defer client.Stop()

	labels := model.LabelSet{
		"level": "info",
	}
	err = client.Handle(labels, time.Now(), "This is a log message")
	if err != nil {
		log.Printf("Failed to send log: %v", err)
	}
}
```

### Configuration

The `Config` struct allows you to customize the client's behavior:

```go
cfg := client.Config{
	URL:       "http://localhost:3100/loki/api/v1/push",
	BatchWait: 1 * time.Second,
	BatchSize: 1024 * 1024, // 1 MB
	Timeout:   10 * time.Second,
	BackoffConfig: client.BackoffConfig{
		MinBackoff: 100 * time.Millisecond,
		MaxBackoff: 10 * time.Second,
		MaxRetries: 5,
	},
	ExternalLabels: model.LabelSet{
		"environment": "production",
	},
}
client, err := client.New(cfg, logger)
```

### Sending Logs

Use the `Handle` method to send logs. Each log entry includes a set of labels, a timestamp, and a log message.

```go
labels := model.LabelSet{
	"level": "error",
	"job":   "worker",
}
err := client.Handle(labels, time.Now(), "An error occurred")
if err != nil {
	log.Printf("Failed to send log: %v", err)
}
```

### Stopping the Client

Call the `Stop` method to gracefully shut down the client and flush any remaining logs.

```go
client.Stop()
```

## API Reference

### `client.Config`

- `URL`: The Loki Push API endpoint.
- `BatchWait`: Maximum time to wait before sending a batch of logs.
- `BatchSize`: Maximum size of a batch in bytes.
- `Timeout`: HTTP request timeout.
- `BackoffConfig`: Configuration for retrying failed requests.
- `ExternalLabels`: Labels to include with every log entry.

### `client.BackoffConfig`

- `MinBackoff`: Initial backoff duration.
- `MaxBackoff`: Maximum backoff duration.
- `MaxRetries`: Maximum number of retries (0 for infinite retries).

### `client.Client`

- `New(cfg Config, logger Logger)`: Creates a new client with the specified configuration.
- `NewWithDefaults(url string, externalLabels model.LabelSet, logger Logger)`: Creates a new client with default configuration.
- `Handle(labels model.LabelSet, timestamp time.Time, line string)`: Adds a log entry to the next batch.
- `Stop()`: Stops the client and flushes any remaining logs.

## Development

### Generating Protobuf Code

The `logproto` package defines the protobuf messages and services used by the client. To regenerate the Go code from the `.proto` file, use the following command:

```sh
protoc --gogofaster_out=. --proto_path=logproto logproto.proto
```

## Protobuf Code Generation

To generate the Go code for the protobuf definitions, ensure you have the `protoc-gen-gogofaster` plugin installed. You can install it using:

```sh
go install github.com/gogo/protobuf/protoc-gen-gogofaster@latest
```

Then, run the following command to generate the protobuf code:

```sh
protoc --gogofaster_out=. --proto_path=logproto logproto.proto
```

## Backoff Configuration Details

The backoff mechanism uses the "Full Jitter" approach for exponential backoff. This means the wait time is randomized between 0 and the current backoff duration, which helps to avoid synchronized retries in distributed systems.

## Error Logging

The `Logger` function is used to log errors and warnings. You can customize this function to integrate with your preferred logging framework. For example:

```go
logger := func(v ...interface{}) {
    log.Printf("[Loki Client] %v", v...)
}
```

### Running Tests

To run the tests, use the following command:

```sh
go test ./...
```

## License

This project is licensed under the [Apache License 2.0](LICENSE).