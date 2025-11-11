package main

import (
	"context"
	"time"

	"github.com/Gustik/metrics/internal/agent"
	"github.com/Gustik/metrics/internal/agent/collector"
	"github.com/Gustik/metrics/internal/agent/sender"
)

func main() {
	collector := collector.NewCollector()
	sender := sender.NewHTTPSender("http://localhost:8080")
	agent := agent.NewAgent(collector, sender)

	agent.Run(context.Background(), 2*time.Second, 10*time.Second)
}
