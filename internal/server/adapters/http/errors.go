package httpserver

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func(wa *WebApp) BindError(ctx *gin.Context, err error) {

	wa.log.Sugar().Errorf("request failed: %s", err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(),})
}