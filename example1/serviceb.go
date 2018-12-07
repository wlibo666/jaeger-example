package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func StartService2() error {
	// tracer配置
	cfg := config.Configuration{
		ServiceName: "service_b",
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
	s1.HandleFunc("/b", func(writer http.ResponseWriter, request *http.Request) {
		var sp opentracing.Span
		context, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
		// 请求头未带trace信息,生成新span
		if err != nil {
			fmt.Printf("Extract failed: %v", err)
			sp = opentracing.StartSpan("/b")
		} else {
			// 携带trace信息生成子span
			sp = opentracing.StartSpan("/b", opentracing.ChildOf(context), opentracing.Tags{"name": "service_b"})
			defer sp.Finish()
		}
		// 异步请求服务c和d
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			// 请求服务c
			defer wg.Done()
			req, _ := http.NewRequest("GET", "http://127.0.0.1:11002/c", nil)
			err := sp.Tracer().Inject(sp.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header))
			if err != nil {
				fmt.Printf("s1 Could not inject span context into header: %v", err)
				writer.Write([]byte(err.Error()))
				return
			}

			_, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("s1 DefaultClient do failed: %v", err)
				writer.Write([]byte(err.Error()))
				return
			}
		}()
		go func() {
			// 请求服务d
			defer wg.Done()
			req, _ := http.NewRequest("GET", "http://127.0.0.1:11003/d", nil)
			if sp != nil {
				err := sp.Tracer().Inject(sp.Context(),
					opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(req.Header))
				if err != nil {
					fmt.Printf("s1 Could not inject span context into header: %v", err)
					writer.Write([]byte(err.Error()))
					return
				}
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("s1 DefaultClient do failed: %v", err)
				writer.Write([]byte(err.Error()))
				return
			}
		}()
		wg.Wait()

		writer.Write([]byte("service_b ok"))
	})
	return http.ListenAndServe("0.0.0.0:11001", s1)
}

func main() {
	StartService2()
}
