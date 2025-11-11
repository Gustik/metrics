package main

import (
	"net/http"

	"github.com/Gustik/metrics/internal/config"
	"github.com/Gustik/metrics/internal/server/handler"
	"github.com/Gustik/metrics/internal/server/repository"
	"github.com/Gustik/metrics/internal/server/service"
)

func main() {
	cfg := config.Load()
	repo := repository.NewInMemoryMetricRepository()
	svc := service.NewMetricService(repo)
	h := handler.NewMetricHandler(svc)

	http.ListenAndServe(cfg.ServerAddress, handler.SetupRouter(h))
}
