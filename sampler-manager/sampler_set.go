package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlibo666/common-lib/webutils"
)

func samplerSetHandler(ctx *gin.Context) {

}

var samplerSetController *webutils.ControllerFunc = &webutils.ControllerFunc{
	ApiType:    API_NEED_SIGN,
	Method:     http.MethodPut,
	HandlePath: "/sampler/:app/:namespace",
	Handler:    samplerGetHandler,
}
