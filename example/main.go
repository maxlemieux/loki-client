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
		"app": "example-app",
	}

	client, err := client.NewWithDefaults("http://localhost:3100/loki/api/v1/push", externalLabels, logger)
	if err != nil {
		log.Fatalf("Failed to create Loki client: %v", err)
	}
	defer client.Stop()

	labels := model.LabelSet{
		"level": "info",
	}
	err = client.Handle(labels, time.Now(), "This is a log message from the example implementation")
	if err != nil {
		log.Printf("Failed to send log: %v", err)
	}
}
