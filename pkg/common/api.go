package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ParamsErrorMsg = "参数错误"
)

type EmptyData struct {
	Code int
}

type ApiData struct {
	Code int
	Data any
	Msg  string
}

func RespEmpty(c *gin.Context) {
	c.JSON(http.StatusOK, EmptyData{})
}

func RespErrorParams(c *gin.Context) {
	c.JSON(http.StatusOK, ApiData{Code: http.StatusBadRequest, Msg: ParamsErrorMsg})
}

func RespError(c *gin.Context, err error) {
	RespErrorCode(c, -1, err)
}

func RespErrorCode(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, ApiData{Code: code, Msg: err.Error()})
}
