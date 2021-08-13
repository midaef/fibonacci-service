FROM golang:alpine as builder

ADD ./ /go/src/fibonacci-service
WORKDIR /go/src/fibonacci-service

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /fibonacci-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /fibonacci-service ./fibonacci-service

RUN mkdir ./config
COPY ./config/prod-config.yaml ./config

EXPOSE 7001

ENTRYPOINT ["./fibonacci-service"]
