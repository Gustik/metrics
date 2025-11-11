package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Gustik/metrics/internal/agent/collector"
	"github.com/Gustik/metrics/internal/agent/sender"
	"github.com/Gustik/metrics/internal/model"
)

type Agent struct {
	mu        sync.RWMutex
	collector collector.MetricCollector
	sender    sender.Sender
	storage   []model.Metric
}

func NewAgent(collector collector.MetricCollector, sender sender.Sender) *Agent {
	return &Agent{
		collector: collector,
		sender:    sender,
		storage:   make([]model.Metric, 0),
	}
}

func (a *Agent) Run(ctx context.Context, pollIntervar, reportInterval time.Duration) {
	pollTicker := time.NewTicker(pollIntervar)
	reportTicker := time.NewTicker(reportInterval)

	for {
		select {
		case <-pollTicker.C:
			a.Poll()
		case <-reportTicker.C:
			a.SendReport()
		case <-ctx.Done():
			return
		}
	}
}

func (a *Agent) Poll() {
	fmt.Print(".")
	metrics := a.collector.Collect()

	a.mu.Lock()
	a.storage = append(a.storage, metrics...)
	a.mu.Unlock()
}

func (a *Agent) SendReport() {
	a.mu.RLock()
	metrics := make([]model.Metric, len(a.storage))
	metrics = append(metrics, a.storage...)
	a.mu.RUnlock()

	a.sender.Send(metrics)
	fmt.Println("sent")
}
