package httpserver

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)


func(wa *WebApp) BodyMeasureMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context){
		pr, pw := io.Pipe()
		
		go func (){
			var totalSize int
			defer pr.Close()

			b := make([]byte, 1024)

			for {
				chunkSize, err := pr.Read(b)
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
		
		tee := io.TeeReader(ctx.Request.Body, pw)
		data, _:= io.ReadAll(tee)
		pw.Close()

		ctx.Request.Body = io.NopCloser(bytes.NewReader(data))
		ctx.Next()
	} 
}




