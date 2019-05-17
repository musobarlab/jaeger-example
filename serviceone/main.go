package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/musobarlab/jaeger-example/helper/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	tracer, closer := tracing.Init("producer", "localhost:5775")
	defer closer.Close()

	http.Handle("/produce", produceHandler(tracer))

	log.Fatal(http.ListenAndServe(":9001", nil))
}

func produceHandler(tracer opentracing.Tracer) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			res.WriteHeader(http.StatusMethodNotAllowed)
			res.Write([]byte("method not allowed"))
			return
		}

		// start tracing
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			log.Printf("error extract span %v\n", err)
		}

		span := tracer.StartSpan("produce", ext.RPCServerOption((spanCtx)))
		defer span.Finish()
		// end tracing

		products := Products{
			Product{ID: "1", Name: "Samsung Galaxy s1"},
			Product{ID: "2", Name: "Samsung J1"},
			Product{ID: "3", Name: "Nokia 6"},
			Product{ID: "4", Name: "IPHONE 6"},
		}

		productsJSON, err := json.Marshal(products)
		if err != nil {
			log.Printf("error marshal product %v\n", err)
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(productsJSON)
	})
}

// Product data
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Products data
type Products []Product
