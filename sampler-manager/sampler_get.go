package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go/thrift-gen/sampling"
	"github.com/wlibo666/common-lib/webutils"
)

const (
	API_NO_SIGN   = 0
	API_NEED_SIGN = 1
)

// jaeger client go http request: http://127.0.0.1:10100/sampler/app/namespace?service=service_a

// return: {"strategyType":"RATE_LIMITING","probabilisticSampling":{"samplingRate":0.1},"rateLimitingSampling":{"maxTracesPerSecond":10}}
func samplerGetHandler(ctx *gin.Context) {
	serviceName, _ := webutils.GetQueryString(ctx, "service")
	cf := &sampling.SamplingStrategyResponse{
		StrategyType: sampling.SamplingStrategyType_RATE_LIMITING,
		ProbabilisticSampling: &sampling.ProbabilisticSamplingStrategy{
			SamplingRate: 0.1,
		},
		RateLimitingSampling: &sampling.RateLimitingSamplingStrategy{
			MaxTracesPerSecond: 10,
		},
	}
	fmt.Fprintf(os.Stdout, "service:%s,sampler strategy:%v\n", serviceName, cf)
	ctx.JSON(http.StatusOK, cf)
}

var samplerGetController *webutils.ControllerFunc = &webutils.ControllerFunc{
	ApiType:    API_NO_SIGN,
	Method:     http.MethodGet,
	HandlePath: "/sampler/:app/:namespace",
	Handler:    samplerGetHandler,
}
