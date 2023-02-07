package app

import (
	"context"
	"fmt"
	"github.com/Nikolay-Yakushev/lychee/config"
	httpserver "github.com/Nikolay-Yakushev/lychee/internal/server/adapters/http"
	"github.com/Nikolay-Yakushev/lychee/pkg/closer"
	"github.com/Nikolay-Yakushev/lychee/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const name string = "app"

type App struct{
	log        *zap.Logger
	conf       *config.Config
	closer     *closer.Closer
	metrics    *metrics.Metrics
	registry   *prometheus.Registry
	descr      string
	closeOrder int
}	

func(a *App) GetDescription() string {
	return a.descr
}

func(a *App) GetCloseOrder() int {
	return a.closeOrder
}

func(a *App) Start() error {
	webapp, err := httpserver.New(
		a.conf,
		a.log,
		a.closer,
		a.registry,
		a.metrics,
	)
		
	if err != nil {
		a.log.Sugar().Errorw("Failed to start webapp", "reason", err)
		err = fmt.Errorf("Failed to start webapp. Reason %w", err)
		return err
	}
	webapp.Start()
	return nil
}

func(a *App) Stop(ctx context.Context) error{
	return nil
}

func New(
	logger   *zap.Logger,
	config   *config.Config,
	closer   *closer.Closer,
	metrics  *metrics.Metrics,
	registry *prometheus.Registry,
	)(*App, error){

	namedLog := logger.Named(name)
	
	a := &App{
		log : namedLog,
		metrics: metrics,
		registry: registry,
		conf: config,
		descr: name,
		closer: closer,
		closeOrder: 2,
	}
	a.closer.Add(a)
	return a, nil
}