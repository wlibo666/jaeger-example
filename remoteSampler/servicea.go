package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func StartService1() error {
	// tracer配置
	cfg := config.Configuration{
		ServiceName: "service_a",
		// 设置为远端采样器
		Sampler: &config.SamplerConfig{
			Type: jaeger.SamplerTypeRemote,
			// 采样URL地址
			SamplingServerURL:       "http://127.0.0.1:10100/sampler/app/namespace",
			SamplingRefreshInterval: 10 * time.Second,
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
	s1.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// 生成一个span
		sp := opentracing.StartSpan("/", opentracing.Tags{"name": "service_a"})
		// span提交到agent
		defer sp.Finish()

		fmt.Printf("sp :%v\n", sp)

		// 将span的Uber-Trace-Id信息注入http头
		req, _ := http.NewRequest("GET", "http://127.0.0.1:11001/b", nil)
		err := sp.Tracer().Inject(sp.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			fmt.Printf("s1 Could not inject span context into header: %v", err)
			writer.Write([]byte(err.Error()))
			return
		}
		// 开始http请求
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("s1 DefaultClient do failed: %v", err)
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write([]byte("service_a ok"))
	})
	return http.ListenAndServe("0.0.0.0:11000", s1)
}

func main() {
	StartService1()
}
