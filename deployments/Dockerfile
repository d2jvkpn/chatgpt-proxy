FROM golang:1 AS builder
LABEL stage=chatgpt-proxy_builder

ARG ldflags=${ldflags}
ENV TZ="Asia/Shanghai"
WORKDIR /opt/chatgpt-proxy

ADD ./main.go ./go.mod ./go.sum ./project.yaml  ./
COPY ./internal ./internal
# COPY ./pkg     ./pkg
# COPY ./vendor  ./vendor

# in alpine, date doesn't parse %:z
RUN go env -w GOPROXY="https://goproxy.cn,direct" && \
  go build -o main -ldflags="-w -s ${ldflags}" main.go

####
FROM registry.cn-shanghai.aliyuncs.com/d2jvkpn/langchain:latest

ENV TZ="Asia/Shanghai"
WORKDIR /home/hello/chatgpt-proxy

COPY --from=builder /opt/chatgpt-proxy/main ./main

EXPOSE 3020
CMD ["./main", "--addr=:3020", "--config=configs/prod.yaml", "--release"]
