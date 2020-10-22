
This is a simple sample of use GO, NodeJs and gRPC to message communication. Both client and server are 
using de file product.proto for gRPC e protobuf definitions.

### GO http server and gRPC client

This sample shows a http server and a gRPC client, it spans 4 workers to make gRPC requests in parallel

+ to run tests you will need GO env installed
```
cd http-product-list
go test ./...
```

### NodeJs

Acting as gRPC server, accept incoming requests for discounts checks for some date related (birthday or blackfriday)

+ to run tests you will need nodejs and npm installed

```
cd grpc-discount
npm run test
```

### 

## Getting Started

For local run, you will need Docker and docker-compose, in bash or prompt do the following steps:

```
cd http-product-list
docker-compose up -d
```

The ports used are 8000 and 50051. 
To test using curl:

```
curl --location --request GET "http://localhost:8000/product"
```