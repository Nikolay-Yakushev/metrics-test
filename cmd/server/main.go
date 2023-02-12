package main

import (
	"context"
	// "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nikolay-Yakushev/lychee/config"
	"github.com/Nikolay-Yakushev/lychee/internal/server/app"
	"github.com/Nikolay-Yakushev/lychee/pkg/closer"
	"github.com/Nikolay-Yakushev/lychee/pkg/logger"
	"github.com/Nikolay-Yakushev/lychee/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const name string = "lychee"


func main() {
	registry := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(registry, name)

	registry.MustRegister(
		collectors.NewGoCollector(),
		metrics.ReqBodySize,
	)
	
	conf := config.New()
	logger, err := logger.New(conf)
	if err != nil {
		logger.Sugar().Fatalf("failed to initialize logger %s", err.Error())
	}

	metrics.ReqBodySize.
			With(prometheus.Labels{"method": "GET"}).
			Observe(float64(12))


	closer := closer.New(logger)
	
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGKILL, os.Interrupt,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT,
	)
	defer cancel()

	logger.Sugar().Infof("Starting %s application", name)

	app, err := app.New(logger, conf, closer, metrics, registry)
	if err != nil {
		logger.Sugar().Fatalw(
			"Failed to initialize app",
			"Reason", err.Error(),
		)
	}

	logger.Sugar().Debug("Start up %s compleated", name)

	e := app.Start()
	if e != nil{
		logger.Sugar().Fatalw(
			"Failed to initialize app",
			"Reason", err.Error(),
		)
	}
	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer stopCancel()

	closer.Close(stopCtx)
}