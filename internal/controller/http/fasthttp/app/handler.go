package app

import (
	"mime"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type Controller struct {
	filePath string
	logger   zerolog.Logger
}

func NewController(filePath string, logger zerolog.Logger) *Controller {
	return &Controller{
		filePath: filePath,
		logger:   logger,
	}
}

func (c *Controller) GetFile(ctx *fasthttp.RequestCtx) {
	fileInfo, err := os.Stat(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		} else {
			ctx.Error("Internal Error", fasthttp.StatusInternalServerError)
		}
		return
	}

	if fileInfo.IsDir() {
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
		return
	}

	ext := filepath.Ext(c.filePath)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		ctx.SetContentType(mimeType)
	} else {
		ctx.SetContentType("application/octet-stream")
	}

	ctx.SendFile(c.filePath)
}
