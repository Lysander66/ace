package controller

import (
	"net/http"

	"github.com/Lysander66/ace/pkg/common"
	"github.com/gin-gonic/gin"
)

// -------------------------*------------------------- crypto -------------------------#-------------------------

func md5Sum(c *gin.Context) {
	SucceedResp(c, common.MD5Sum(c.Query("s")))
}

func sha1Sum(c *gin.Context) {
	SucceedResp(c, common.SHA1Sum(c.Query("s")))
}

func sha256Sum(c *gin.Context) {
	SucceedResp(c, common.SHA256Sum(c.Query("s")))
}

func sha512Sum(c *gin.Context) {
	SucceedResp(c, common.SHA512Sum(c.Query("s")))
}

// -------------------------*------------------------- response -------------------------#-------------------------

type ApiData struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type ApiListData struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

func SucceedResp(c *gin.Context, data any) {
	c.JSON(http.StatusOK, ApiData{Data: data})
}

func ErrorResp(c *gin.Context, err error) {
	ErrorCodeResp(c, -1, err)
}

func ErrorCodeResp(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, ApiData{Code: code, Msg: err.Error()})
}
