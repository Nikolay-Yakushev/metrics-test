package httpserver

import (
	"io"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)


func(wa *WebApp) BodyMeasureMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context){
		pr, pw := io.Pipe()
		tee := io.TeeReader(ctx.Request.Body, pw)

		go func (){
			defer pw.Close()
			var totalSize int
			b := make([]byte, 1024)

			for {
				chunkSize, err := tee.Read(b)
				if err == io.EOF{
					break
				}
				if err == io.ErrUnexpectedEOF{
					totalSize+=chunkSize
					break
				}
				totalSize+=chunkSize
			}
			wa.log.Sugar().Infof("total size %d", totalSize)
			wa.metrics.
			ReqBodySize.With(
				prometheus.Labels{"method": ctx.Request.Method}).
				Observe(float64(totalSize))				
		}()

		ctx.Request.Body = io.NopCloser(pr)
		ctx.Next()
	} 
}




