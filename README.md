### Distributed Tracing Example with Jaeger

- https://opentracing.io/
- https://www.jaegertracing.io/

#### Getting started

Clone this project
```shell
$ git clone https://github.com/musobarlab/jaeger-example.git
$ go get -u
```

Start Jaeger Agent, Jaeger Collector, Jaeger Query and Storage (we will use Elastic Search) 
```shell
$ docker-compose up -d elasticsearch
$ docker-compose up
```

open Jaeger UI http://127.0.0.1:16686

Start our two Microservices server example

- microservice one

```shell
$ cd serviceone
$ go run main.go
```

- microservice two

```shell
$ cd servicetwo
$ go run main.go
```

Run our client app, in this example we will search product with name `Nokia`

```shell
$ cd client
$ go run main.go Nokia
```

Refresh the Jaeger UI http://127.0.0.1:16686, and you will see the tracing output