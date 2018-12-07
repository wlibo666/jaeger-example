package main

import (
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func StartService3() error {
	// tracer配置
	cfg := config.Configuration{
		ServiceName: "service_c",
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// report配置信息,包括agent地址
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	// 设置全局tracer
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Printf("new tracer failed,err:%s\n", err.Error())
		return err
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	s1 := http.NewServeMux()
	s1.HandleFunc("/c", func(writer http.ResponseWriter, request *http.Request) {
		var sp opentracing.Span
		context, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
		// 请求头未带trace信息,生成新span
		if err != nil {
			fmt.Printf("Extract failed: %v", err)
			sp = opentracing.StartSpan("/c")
		} else {
			// 携带trace信息生成子span
			sp = opentracing.StartSpan("/c", opentracing.ChildOf(context), opentracing.Tags{"name": "service_c"})
			defer sp.Finish()
		}

		writer.Write([]byte("service_c ok"))
	})
	return http.ListenAndServe("0.0.0.0:11002", s1)
}

func main() {
	StartService3()
}
