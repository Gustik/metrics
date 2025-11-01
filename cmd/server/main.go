package main

import (
	"net/http"

	"github.com/Gustik/metrics/internal/config"
	"github.com/Gustik/metrics/internal/handler"
	"github.com/Gustik/metrics/internal/repository"
	"github.com/Gustik/metrics/internal/service"
)

func main() {
	cfg := config.Load()
	repo := repository.NewInMemoryMetricRepository()
	svc := service.NewMetricService(repo)
	h := handler.NewMetricHandler(svc)

	http.ListenAndServe(cfg.ServerAddress, handler.SetupRouter(h))
}
