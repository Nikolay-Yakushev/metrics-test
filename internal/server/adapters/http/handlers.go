package httpserver

import (
	"io"
	"mime"
	"net/http"
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

func(wa *WebApp) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"health": true})
}


func(wa *WebApp) upload (ctx *gin.Context){
	multipart, err := ctx.Request.MultipartReader()

	if err != nil{
		wa.log.Sugar().Fatalw("Failed to create MultipartReader", "Error: ", err)
	}
	for {
		mimePart, err := multipart.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			wa.log.Sugar().Errorw("Error reading multipart section", "Error: ", err)
			break
		}
		_, params, err := mime.ParseMediaType(mimePart.Header.Get("Content-Disposition"))
		if err != nil {
			wa.log.Sugar().Errorw("Invalid Content-Disposition", "Error: ", err)
			break
		}

		// TODO how to make writing multipart in gorotines?
		curdir, err := os.Getwd()
		if err != nil {
			wa.log.Sugar().Errorf("Failed to get current working directory", "Error: ", err)
		}

		filename := filepath.Join(curdir, params["filename"])
		file, err := os.Create(filename)
		defer file.Close()

		if err != nil {
			wa.log.Sugar().Errorf("Failed to create file", filename, "Error: ", err)
		}
		data, err := ioutil.ReadAll(mimePart)
		if err != nil {
			wa.log.Sugar().Errorf("Failed to create file", filename, "Error: ", err)
		}
		file.Write(data)
		// TODO how to make sure all gorotines have finished succsessfully before sending replying to user?
	}

	ctx.JSON(http.StatusOK, gin.H{"uploaded": true})
}
	