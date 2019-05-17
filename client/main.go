package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/musobarlab/jaeger-example/helper"
	"github.com/musobarlab/jaeger-example/helper/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	tracer, closer := tracing.Init("client-service", "localhost:5775")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	productName := os.Args[1]

	span := tracer.StartSpan("search-product")
	span.SetTag("search-to", productName)
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// simple example
	//helloStr := formatString(ctx, helloTo)
	//printHello(ctx, helloStr)

	//http request example
	data := produce(ctx)
	product := consume(ctx, data, productName)
	fmt.Println("product id : ", product.ID)
	fmt.Println("product name : ", product.Name)

}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}

func produce(ctx context.Context) []byte {
	span, _ := opentracing.StartSpanFromContext(ctx, "produce")
	defer span.Finish()

	url := "http://localhost:9001/produce"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	// send span to server
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := helper.Do(req)
	if err != nil {
		panic(err.Error())
	}

	span.LogFields(
		log.String("event", "produce"),
		log.String("url", url),
		log.String("value", string(resp)),
	)

	return resp
}

func consume(ctx context.Context, data []byte, search string) Product {
	span, _ := opentracing.StartSpanFromContext(ctx, "consume")
	defer span.Finish()

	url := fmt.Sprintf("http://localhost:9002/consume?search=%s", search)
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err.Error())
	}

	// send span to server
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "POST")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := helper.Do(req)
	if err != nil {
		panic(err.Error())
	}

	span.LogFields(
		log.String("event", "consume"),
		log.String("url", url),
		log.String("search", search),
		log.String("value", string(resp)),
	)

	var product Product
	err = json.Unmarshal(resp, &product)

	return product
}

// Product data
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
