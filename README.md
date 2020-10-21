# go-node-grpc-product

## Getting Started

* Run docker compose

```
cd http-product-list
docker-compose up -d
```

The ports used are 8000 and 50051. 
To test using curl:

```
curl --location --request GET "http://localhost:8000/product"
```