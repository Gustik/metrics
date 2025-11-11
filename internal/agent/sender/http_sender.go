package sender

import (
	"fmt"
	"net/http"

	"github.com/Gustik/metrics/internal/model"
)

type Sender interface {
	Send(metrics []model.Metric)
}

type sender struct {
	url string
}

func NewHTTPSender(url string) Sender {
	return &sender{url: url}
}

func (s sender) Send(metrics []model.Metric) {
	for _, metric := range metrics {
		if metric.MType == model.Counter {
			u := fmt.Sprintf("%s/update/counter/%s/%d", s.url, metric.ID, *metric.Delta)
			fmt.Println(u)
			http.Post(u, "Content-Type: text/plain", nil)
		}

		if metric.MType == model.Gauge {
			u := fmt.Sprintf("%s/update/gauge/%s/%f", s.url, metric.ID, *metric.Value)
			http.Post(u, "Content-Type: text/plain", nil)
		}
	}
}
