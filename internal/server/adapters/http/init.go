package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Nikolay-Yakushev/lychee/config"
	"github.com/Nikolay-Yakushev/lychee/pkg/closer"
	"github.com/Nikolay-Yakushev/lychee/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const name string = "webapp"

type WebApp struct {
	registry   *prometheus.Registry
	server     *http.Server
	log        *zap.Logger
	cfg        *config.Config
	closer     *closer.Closer
	metrics    *metrics.Metrics
	listener   net.Listener
	once       sync.Once
	wg         sync.WaitGroup
	closeOrder int
	descr      string
}


func(wa *WebApp) GetDescription() string {
	return wa.descr
}


func(wa *WebApp) GetCloseOrder() int {
	return wa.closeOrder
}


func(wa *WebApp) Start() error {
	var err error

	go func() {
		err = wa.server.Serve(wa.listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) { 
			wa.log.Sugar().Fatalw("Server startup failed", "reason", err)
		}
	}()
	return err
}

func (wa *WebApp) Stop(ctx context.Context) error {
	var err error

	wa.once.Do(func() {
		err = wa.server.Shutdown(ctx)
	})

	if err != nil{
		wa.log.Sugar().Errorw("failed to stop %s", wa.GetDescription(), "Reason", err)
		err := fmt.Errorf("Failed to stop %s. Reason: %w", wa.GetDescription(), err)
		return err
	}
	return nil
}

func New(
	config   *config.Config,
	logger   *zap.Logger,
	closer   *closer.Closer,
	registry *prometheus.Registry,
	metrics  *metrics.Metrics,
	)(*WebApp, error){
	

	stringPort := strconv.Itoa(config.ServerPort)

	listener, err := net.Listen("tcp", ":" + stringPort)
	if err != nil{
		logger.Sugar().Errorw("Failed to initialize WebApp", "error", err)
		return nil, fmt.Errorf("WebApp init failed: %w", err)
	}
	router := gin.New()

	server := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort),
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	
	namedLog := logger.Named(name)
	wa := &WebApp{
		server: server,
		log : namedLog,
		cfg: config,
		descr: name,
		listener: listener,
		closer: closer,
		closeOrder: 1,
		registry: registry,
		metrics: metrics,
	}
	wa.initRoutes(router, logger)
	closer.Add(wa)
	return wa, nil
}