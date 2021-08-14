# FIBONACCI-SERVICE

## PORT: 7001

## LOCAL CONFIGURATION

### REQUIREMENTS
- go 1.16
- docker & docker-compose
- redis
- evans or bloom rpc
- postman or curl

## DOCKER

### Build
```shell
docker-compose build
```
### Run
```shell
docker-compose up
```

## Go

### Install [REDIS](https://redis.io/download)

Create **local-config.yaml** file in config directory:
```yaml
appconfig:
  ip: "0.0.0.0"
  grpc_port: "7001"
  http_port: "8080"
redis:
  db: 0
  host: "0.0.0.0"
  port: "6379"
  password: ""
```

### Build
```shell
go build -o fibonacci-service cmd/main.go -config-path ./config/local-config.yaml
```
### Run
```shell
./fibonacci-service
```

### OR

```shell
go run cmd/main.go -config-path ./config/local-config.yaml
```

## Testing GRPC

### Install [EVANS](https://github.com/ktr0731/evans) or [BloomRPC](https://github.com/uw-labs/bloomrpc)

### Run for evans with Makefile
```shell
make evans name=fibonacci port=7001

call FibonacciSequences
```

### OR

### Run for evans without Makefile
```shell
evans api/fibonacci.proto -p 7001

call FibonacciSequences
```

## Testing HTTP

```shell
POST http://localhost:8080/api/v1/fibonacci_sequences
```

### Body

```shell
{
  "x":5,
  "y":15
}
```

### Install [POSTMAN](https://www.postman.com/downloads/)

### OR

### Run curl 
```shell
curl -XPOST -H "Content-type: application/json" -d '{"x":5, "y":15}' 'http://localhost:8080/api/v1/fibonacci_sequences'
```