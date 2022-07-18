package router

import (
	"bytes"
	"fmt"
	"imgconv/errors"
	"imgconv/manager"
	"imgconv/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	manager manager.ConversionManager
}

func (r *Router) reportError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, &ErrorResponse{
		Error: err.Error(),
	})
}

func (r *Router) OnHEIFVersionRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": utils.GetLibHeifVersion(),
	})
}

func (r *Router) OnSingleFileConversionRequest(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		r.reportError(c, err)
		return
	}
	f, err := file.Open()
	if err != nil {
		r.reportError(c, err)
		return
	}
	var inp bytes.Buffer
	_, err = io.Copy(&inp, f)
	if err != nil {
		r.reportError(c, err)
		return
	}

	id, err := r.manager.AddConversion(&manager.AddConversionOptions{
		InputBytes: inp.Bytes(),
		Format:     "PNG",
		Filename:   file.Filename,
	})

	status := r.manager.GetConversionStatusById(id)
	c.JSON(http.StatusCreated, &ConversionStatusResponse{
		Status:   status.Status,
		Filename: status.Filename,
		Format:   status.Format,
		Id:       status.Id,
	})
}

func (r *Router) OnGetConversionStatusRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.reportError(c, errors.IdParamNotFound)
		return
	}
	status := r.manager.GetConversionStatusById(id)
	if status == nil {
		c.JSON(http.StatusNotFound, &ErrorResponse{
			Error: errors.ConversionTaskNotFound.Error(),
		})
		return
	}
	resp := &ConversionStatusResponse{
		Status:   status.Status,
		Filename: status.Filename,
		Format:   status.Format,
		Id:       status.Id,
	}
	if status.Error != nil {
		resp.Error = status.Error.Error()
	}
	c.JSON(http.StatusFound, resp)
}

func (r *Router) OnDownloadConvertedImageRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		r.reportError(c, errors.IdParamNotFound)
		return
	}
	status := r.manager.GetConversionStatusById(id)
	if status == nil {
		c.JSON(http.StatusNotFound, &ErrorResponse{
			Error: errors.ConversionTaskNotFound.Error(),
		})
		return
	}
	switch status.Status {
	case manager.TypeStatusProcessing:
		c.JSON(http.StatusTooEarly, &ConversionStatusResponse{
			Status:   status.Status,
			Filename: status.Filename,
			Format:   status.Format,
			Id:       status.Id,
		})
		break
	case manager.TypeStatusFailed:
		c.JSON(http.StatusBadRequest, &ConversionStatusResponse{
			Status:   status.Status,
			Filename: status.Filename,
			Format:   status.Format,
			Id:       status.Id,
			Error:    status.Error.Error(),
		})
		break
	case manager.TypeStatusCompleted:
		outBytes, err := r.manager.GetOutputImageBytesById(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, &ConversionStatusResponse{
				Status:   status.Status,
				Filename: status.Filename,
				Format:   status.Format,
				Id:       status.Id,
				Error:    status.Error.Error(),
			})
			break
		}
		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, status.Filename),
		}
		c.DataFromReader(http.StatusOK, int64(len(outBytes)), "image/png", bytes.NewReader(outBytes), extraHeaders)
		break
	}
}

func NewRouter(manager manager.ConversionManager) *Router {
	return &Router{
		manager: manager,
	}
}
