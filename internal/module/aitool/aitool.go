package aitool

import (
	"CloudPhoto/internal/module/storage"
	"CloudPhoto/internal/tool"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

func ChangeFace(c *gin.Context) {
	face, err := c.FormFile("face")
	tool.HandleErr(err, func() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
		return
	})

	bodyId := c.Query("body_id")
	if bodyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing base_id"})
		return
	}
	doneCh := make(chan string)
	errCh := make(chan error)
	go func() {
		bodyFile, err := os.Open(storage.GetBodyFilePath(bodyId))
		tool.HandleErr(err, func() { errCh <- err; return })
		defer tool.PanicIfErr(bodyFile.Close())
		bodyByte, err := io.ReadAll(bodyFile)
		tool.HandleErr(err, func() { errCh <- err; return })
		faceFile, err := face.Open()
		tool.HandleErr(err, func() { errCh <- err; return })
		defer tool.PanicIfErr(faceFile.Close())
		faceByte, err := io.ReadAll(faceFile)
		tool.HandleErr(err, func() { errCh <- err; return })
		callPhotosAiApi("changeFaceApi", doneCh, errCh, &bodyByte, &faceByte)
	}()
	select {
	case resultPath := <-doneCh:
		c.File(resultPath)
	case err := <-errCh:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "AI processing timeout"})
	}
}

func AddFigure(c *gin.Context) {
	files, err := c.MultipartForm()
	tool.HandleErr(err, func() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing files"})
		return
	})
	backGround, ok := files.File["backGround"]
	if !ok || len(backGround) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing backGround"})
		return
	} else if len(backGround) > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many backGround"})
		return
	}
	figures, ok := files.File["figures"]
	if !ok || len(figures) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing figures"})
		return
	} else if len(figures) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many figures"})
	}
	doneCh := make(chan string)
	errCh := make(chan error)
	go func() {
		backGroundFile, err := backGround[0].Open()
		tool.HandleErr(err, func() { errCh <- err; return })
		defer tool.PanicIfErr(backGroundFile.Close())
		backGroundByte, err := io.ReadAll(backGroundFile)
		figuresByte := make([]*[]byte, len(figures))
		for i, figure := range figures {
			figureFile, err := figure.Open()
			tool.HandleErr(err, func() { errCh <- err; return })
			defer tool.PanicIfErr(figureFile.Close())
			*figuresByte[i], err = io.ReadAll(figureFile)
		}

		callPhotosAiApi("addFiguresApi", doneCh, errCh, &backGroundByte, figuresByte...)
	}()
	select {
	case resultPath := <-doneCh:
		c.File(resultPath)
	case err := <-errCh:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "AI processing timeout"})
	}
}

func callPhotosAiApi(api string, doneCh chan<- string, errCh chan<- error, photo *[]byte, photos ...*[]byte) {
	defer close(doneCh)
	defer close(errCh)
	resultPath, err :=
		//接api
		func(f *[]byte, fs ...*[]byte) (string, error) {
			return "", nil
		}(photo, photos...)
	tool.HandleErr(err, func() { errCh <- err; return })
	doneCh <- resultPath
}
