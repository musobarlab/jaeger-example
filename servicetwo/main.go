package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/musobarlab/jaeger-example/helper/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	tracer, closer := tracing.Init("consumer", "localhost:5775")
	defer closer.Close()

	http.Handle("/consume", produceHandler(tracer))

	log.Fatal(http.ListenAndServe(":9002", nil))
}

func produceHandler(tracer opentracing.Tracer) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		if req.Method != "POST" {
			res.WriteHeader(http.StatusMethodNotAllowed)
			res.Write([]byte("method not allowed"))
			return
		}

		// start tracing
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			log.Printf("error extract span %v\n", err)
		}

		span := tracer.StartSpan("consume", ext.RPCServerOption((spanCtx)))
		defer span.Finish()
		// end tracing

		search, ok := req.URL.Query()["search"]

		if !ok || len(search) < 1 {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("search param required"))
			return
		}

		var (
			products Products
			product  Product
		)
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&products)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("error encode body"))
			return
		}

		// find product by search
		for _, v := range products {
			if strings.Contains(v.Name, search[0]) {
				product = v
				break
			}
		}

		productJSON, err := json.Marshal(product)
		if err != nil {
			log.Printf("error marshal product %v\n", err)
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(productJSON)
	})
}

// Product data
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Products data
type Products []Product
