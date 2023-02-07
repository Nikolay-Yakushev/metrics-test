package httpserver

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)




func(wa *WebApp) initRoutes(router *gin.Engine, logger *zap.Logger){
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(wa.BodyMeasureMiddleware())
	
	router.GET( "/metrics", gin.WrapH(promhttp.HandlerFor(
		wa.registry, promhttp.HandlerOpts{Registry: wa.registry})))

	v1 := router.Group("/api/v1/")

	{
		v1.GET("/health", wa.health)
		v1.POST("/upload", wa.upload)
	}

}